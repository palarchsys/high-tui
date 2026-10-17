#!/usr/bin/env python3
"""Read theme.css and the window JSON files, then write layout.tsv."""

import json
import os
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent
KNOWN = {
    "color",
    "background",
    "background-color",
    "border",
    "border-width",
    "border-style",
    "border-color",
    "border-radius",
    "padding",
    "margin",
    "max-width",
    "column-gap",
    "line-height",
    "font-weight",
}
TAGS = {"window", "header", "nav", "select", "input", "textarea", "button", "command", "status", "main"}
ACTIONS = {"load", "edit", "window", "return", "save", "refresh", "run", "set"}
ACTION_KEYS = {"action", "target", "file", "name", "href", "run", "output", "value"}
COL_KEYS = {"id", "tag", "fr", "label", "class", "file", "name", "save", "readonly", "focus", "enter", "items"}
ITEM_KEYS = {"id", "label", "value", "focus", "enter"}
WIN_KEYS = {"id", "tag", "label", "class", "rows"}
ROW_KEYS = {"id", "tag", "class", "lines", "columns"}


def die(msg):
    print(f"erreur : {msg}", file=sys.stderr)
    raise SystemExit(1)


def strip_comments(text):
    out = []
    i = 0
    while i < len(text):
        if text.startswith("/*", i):
            j = text.find("*/", i + 2)
            if j < 0:
                die("commentaire CSS non fermé")
            i = j + 2
            continue
        out.append(text[i])
        i += 1
    return "".join(out)


def length_of(token, where):
    raw = token.strip()
    if raw.endswith("ch"):
        raw = raw[:-2]
    if not raw.isdigit():
        die(f"longueur invalide « {token.strip()} » ({where})")
    return int(raw)


def four_edges(parts, where):
    nums = [length_of(p, where) for p in parts]
    if len(nums) == 1:
        return (nums[0], nums[0], nums[0], nums[0])
    if len(nums) == 2:
        return (nums[0], nums[1], nums[0], nums[1])
    if len(nums) == 3:
        return (nums[0], nums[1], nums[2], nums[1])
    if len(nums) == 4:
        return (nums[0], nums[1], nums[2], nums[3])
    die(f"trop de valeurs ({where})")


def parse_selector(part, where):
    text = part.strip()
    if not text:
        die(f"sélecteur vide ({where})")
    tag = cls = eid = ""
    focus = False
    i = 0
    n = len(text)
    if i < n and (text[i].isalpha() or text[i] == "_"):
        j = i + 1
        while j < n and (text[j].isalnum() or text[j] in "-_"):
            j += 1
        tag = text[i:j]
        i = j
    if i < n and text[i] == ".":
        j = i + 1
        while j < n and (text[j].isalnum() or text[j] in "-_"):
            j += 1
        if j == i + 1:
            die(f"classe vide ({where})")
        cls = text[i + 1:j]
        i = j
    if i < n and text[i] == "#":
        j = i + 1
        while j < n and (text[j].isalnum() or text[j] in "-_"):
            j += 1
        if j == i + 1:
            die(f"identifiant CSS vide ({where})")
        eid = text[i + 1:j]
        i = j
    if text.startswith(":focus", i):
        focus = True
        i += 6
    if i != n or not (tag or cls or eid or focus):
        die(f"sélecteur refusé « {text} » ({where})")
    if tag and tag not in TAGS:
        die(f"balise inconnue « {tag} » ({where})")
    spec = ((1 if eid else 0), (1 if cls else 0) + (1 if focus else 0), 1 if tag else 0)
    return {"tag": tag, "class": cls, "id": eid, "focus": focus, "spec": spec}


def parse_decls(body, where):
    props = []
    for piece in body.split(";"):
        if not piece.strip():
            continue
        if ":" not in piece:
            die(f"déclaration invalide « {piece.strip()} » ({where})")
        name, val = piece.split(":", 1)
        name = name.strip().lower()
        val = " ".join(val.strip().split())
        if name not in KNOWN:
            die(f"propriété inconnue « {name} » ({where})")
        props.append((name, val))
    return props


def parse_css(text):
    raw = strip_comments(text)
    rules = []
    i = 0
    order = 0
    while True:
        while i < len(raw) and raw[i].isspace():
            i += 1
        if i >= len(raw):
            break
        j = raw.find("{", i)
        if j < 0:
            die("accolade ouvrante manquante dans le CSS")
        sel_text = raw[i:j].strip()
        k = raw.find("}", j)
        if k < 0:
            die("accolade fermante manquante dans le CSS")
        body = raw[j + 1:k]
        i = k + 1
        if sel_text.startswith("@"):
            die(f"règle @ refusée ({sel_text.split()[0]})")
        for part in sel_text.split(","):
            sel = parse_selector(part, "theme.css")
            rules.append((sel, parse_decls(body, "theme.css"), order))
            order += 1
    return rules


def blank_style():
    return {
        "color": "",
        "background": "",
        "bw": 0,
        "bc": "",
        "pad": (0, 0, 0, 0),
        "mar": (0, 0, 0, 0),
        "lh": 1,
        "bold": 0,
        "gap": 0,
        "max_w": 260,
    }


def apply_decl(style, name, val, where):
    if name in ("color",):
        style["color"] = val
    elif name in ("background", "background-color"):
        style["background"] = val
    elif name == "border-color":
        style["bc"] = val
    elif name == "border-width":
        style["bw"] = length_of(val, where)
    elif name == "border-style":
        if val not in ("solid", "none"):
            die(f"border-style refusé « {val} » ({where})")
        if val == "none":
            style["bw"] = 0
    elif name == "border-radius":
        length_of(val, where)
    elif name == "border":
        bits = val.split()
        width = color = None
        for bit in bits:
            if bit in ("solid", "none"):
                if bit == "none":
                    width = 0
                continue
            if bit.startswith("#"):
                color = bit
                continue
            width = length_of(bit, where)
        if width is None or color is None:
            die(f"border incomplet « {val} » ({where})")
        style["bw"] = width
        style["bc"] = color
    elif name == "padding":
        style["pad"] = four_edges(val.split(), where)
    elif name == "margin":
        style["mar"] = four_edges(val.split(), where)
    elif name == "max-width":
        style["max_w"] = length_of(val, where)
    elif name == "column-gap":
        style["gap"] = length_of(val, where)
    elif name == "line-height":
        if not val.isdigit() or int(val) < 1:
            die(f"line-height refusé « {val} » ({where})")
        style["lh"] = int(val)
    elif name == "font-weight":
        if val not in ("bold", "normal"):
            die(f"font-weight refusé « {val} » ({where})")
        style["bold"] = 1 if val == "bold" else 0
    if style["bw"] not in (0, 1):
        die(f"border-width doit valoir 0 ou 1 ({where})")


def matches(sel, tag, classes, eid, focus):
    if sel["focus"] and not focus:
        return False
    if sel["tag"] and sel["tag"] != tag:
        return False
    if sel["id"] and sel["id"] != eid:
        return False
    if sel["class"] and sel["class"] not in classes:
        return False
    return True


def resolve(rules, tag, classes, eid, focus, inherit):
    style = blank_style()
    style["color"] = inherit.get("color", "")
    style["background"] = inherit.get("background", "")
    style["bc"] = inherit.get("color", "")
    chosen = []
    for sel, decls, order in rules:
        if matches(sel, tag, classes, eid, focus):
            chosen.append((sel["spec"], order, decls))
    chosen.sort(key=lambda item: (item[0], item[1]))
    where = eid or tag
    for _spec, _order, decls in chosen:
        for name, val in decls:
            apply_decl(style, name, val, where)
    if not style["color"] or not style["background"]:
        die(f"couleur manquante pour {where}")
    if not style["bc"]:
        style["bc"] = style["color"]
    return style


def outer_h(style, lines):
    between = (lines - 1) * (style["lh"] - 1)
    top, _right, bottom, _left = style["pad"]
    return style["bw"] * 2 + top + lines + between + bottom


def text_w(outer, style):
    _top, right, _bottom, left = style["pad"]
    return outer - style["bw"] * 2 - right - left


def column_widths(inner, frs, gap, where):
    n = len(frs)
    if n < 1:
        die(f"rangée sans colonne ({where})")
    if n == 1:
        if inner < 1:
            die(f"largeur insuffisante ({where})")
        return [inner]
    usable = inner - (n - 1) * gap
    total = sum(frs)
    if total < 1 or usable < n:
        die(f"largeur insuffisante ({where})")
    widths = []
    used = 0
    for fr in frs[:-1]:
        width = (usable * fr) // total
        if width < 1:
            die(f"colonne trop étroite ({where})")
        widths.append(width)
        used += width
    last = usable - used
    if last < 1:
        die(f"dernière colonne trop étroite ({where})")
    widths.append(last)
    return widths


def need(obj, key, where):
    if key not in obj:
        die(f"attribut manquant « {key} » ({where})")
    return obj[key]


def load_json(path):
    try:
        data = json.loads(path.read_text(encoding="utf-8"))
    except json.JSONDecodeError as exc:
        die(f"JSON invalide dans {path.name} : {exc.msg}")
    if not isinstance(data, dict):
        die(f"la fenêtre {path.name} doit être un objet")
    return data


def empty_action():
    return {
        "action": "",
        "target": "",
        "file": "",
        "name": "",
        "href": "",
        "run": "",
        "output": "",
        "value": "",
    }


def reject_unknown(obj, allowed, where):
    if not isinstance(obj, dict):
        die(f"objet attendu ({where})")
    for key in obj:
        if key in allowed:
            continue
        if key in ("title", "text"):
            die(f"utilisez « label » à la place de « {key} » ({where})")
        if key == "action":
            die(f"utilisez « focus » ou « enter » à la place de « action » ({where})")
        die(f"attribut inconnu « {key} » ({where})")


def class_list(raw, where):
    classes = raw.get("class", "")
    if isinstance(classes, str):
        return [classes] if classes else []
    if isinstance(classes, list) and all(isinstance(item, str) for item in classes):
        return classes
    die(f"classe invalide ({where})")


def resolve_ref(text, where):
    if text.count("::") < 2:
        return text
    path_s, name, kind = text.rsplit("::", 2)
    if not path_s or not re.fullmatch(r"[A-Za-z_][A-Za-z0-9_]*", name):
        return text
    if kind not in ("label", "value"):
        die(f"le dernier mot n'est ni label ni value ({where})")
    file_path = Path(path_s)
    if not file_path.is_absolute():
        file_path = ROOT / path_s
    if not file_path.is_file():
        die(f"fichier introuvable : {path_s}")
    found_value = None
    found_label = None
    try:
        raw = file_path.read_text(encoding="utf-8")
    except OSError:
        die(f"fichier introuvable : {path_s}")
    for line in raw.splitlines():
        matched = re.fullmatch(r"([A-Za-z_][A-Za-z0-9_]*)=(.*)", line)
        if matched and matched.group(1) == name and found_value is None:
            found_value = matched.group(2)
        lab = re.fullmatch(rf"# label {re.escape(name)} (.*)", line)
        if lab and found_label is None:
            found_label = lab.group(1)
    if kind == "value":
        if found_value is None:
            die(f"variable manquante ({where})")
        return found_value
    if found_label is None or found_label == "":
        die(f"label manquant ({where})")
    return found_label


def parse_action(raw, where):
    if raw is None:
        return empty_action()
    reject_unknown(raw, ACTION_KEYS, where)
    action = str(raw.get("action", ""))
    act = empty_action()
    act["action"] = action
    for key in ("target", "file", "name", "href", "run", "output", "value"):
        if key in raw and raw[key] is not None:
            act[key] = resolve_ref(str(raw[key]), where) if key == "value" else str(raw[key])
            if "\t" in act[key] or "\n" in act[key]:
                die(f"valeur illégale « {key} » ({where})")
    if not action:
        return act
    if action not in ACTIONS:
        die(f"action inconnue « {action} » ({where})")
    if action == "load":
        if not act["target"]:
            die(f"cible manquante ({where})")
        if not act["file"]:
            die(f"fichier manquant ({where})")
        if not (ROOT / act["file"]).is_file():
            die(f"fichier introuvable : {act['file']}")
    elif action == "edit":
        if not act["target"]:
            die(f"cible manquante ({where})")
    elif action == "window":
        if not act["href"]:
            die(f"href manquant ({where})")
    elif action == "run":
        if not act["run"]:
            die(f"commande manquante ({where})")
    return act


def finish_set(act, fallback, where):
    if act["action"] != "set":
        return act
    if not act["value"]:
        act["value"] = fallback
    if not act["file"]:
        die(f"fichier manquant ({where})")
    if not (ROOT / act["file"]).is_file():
        die(f"fichier introuvable : {act['file']}")
    if not act["name"]:
        die(f"variable manquante ({where})")
    if not act["value"]:
        die(f"valeur manquante ({where})")
    if "\t" in act["value"] or "\n" in act["value"]:
        die(f"valeur illégale ({where})")
    return act


def edit_self(col_id):
    act = empty_action()
    act["action"] = "edit"
    act["target"] = col_id
    return act


def parse_item(raw, ids, col_id):
    reject_unknown(raw, ITEM_KEYS, col_id)
    iid = str(need(raw, "id", col_id))
    if iid in ids:
        die(f"identifiant en double : {iid}")
    ids.add(iid)
    if "label" not in raw:
        die(f"label manquant ({iid})")
    label = resolve_ref(str(raw["label"]), iid)
    value = resolve_ref(str(raw.get("value", "")), iid)
    if label == "" or "\t" in label or "\n" in label:
        die(f"label illégal ({iid})")
    if "\t" in value or "\n" in value:
        die(f"valeur illégale ({iid})")
    return {
        "id": iid,
        "label": label,
        "value": value,
        "focus": finish_set(parse_action(raw.get("focus"), iid), value, iid),
        "enter": finish_set(parse_action(raw.get("enter"), iid), value, iid),
    }


def ensure_return(columns, ids, where):
    lists = [col for col in columns if col["tag"] in ("nav", "select")]
    for col in lists:
        for item in col["items"]:
            if item["focus"]["action"] == "return" or item["enter"]["action"] == "return":
                return
    if not lists:
        die(f"aucune liste pour placer Retour ({where})")
    nav = lists[0]
    item_id = f"{nav['id']}-retour"
    while item_id in ids:
        item_id += "-r"
    ids.add(item_id)
    nav["items"].append({
        "id": item_id,
        "label": "Retour",
        "value": "",
        "focus": empty_action(),
        "enter": {**empty_action(), "action": "return"},
    })


def check_targets(columns):
    zones = {col["id"]: col for col in columns}
    hrefs = []

    def one(act):
        action = act["action"]
        if action in ("load", "edit"):
            target = act["target"]
            if target not in zones:
                die(f"cible absente : {target}")
            tag = zones[target]["tag"]
            if tag not in ("input", "textarea"):
                die(f"cible incompatible : {target}")
            if action == "edit" and tag == "textarea" and zones[target]["readonly"]:
                die(f"cible non éditable : {target}")
        if action == "set" and act["target"]:
            target = act["target"]
            if target not in zones or zones[target]["tag"] not in ("input", "textarea"):
                die(f"cible incompatible : {target}")
        if action == "run" and act["output"]:
            output = act["output"]
            if output not in zones or zones[output]["tag"] != "textarea":
                die(f"cible absente : {output}")
        if action == "window":
            hrefs.append(act["href"])

    for col in columns:
        one(col["focus"])
        one(col["enter"])
        for item in col["items"]:
            one(item["focus"])
            one(item["enter"])
    return hrefs


def parse_window(path, rules, windows, seen_paths):
    path = path.resolve()
    if path in seen_paths:
        return
    seen_paths.add(path)
    data = load_json(path)
    where = path.name
    reject_unknown(data, WIN_KEYS, where)
    if str(need(data, "tag", where)) != "window":
        die(f"balise window attendue ({where})")
    win_id = str(need(data, "id", where))
    if any(w["id"] == win_id for w in windows):
        die(f"identifiant de fenêtre en double : {win_id}")
    ids = {win_id}
    win_style = resolve(rules, "window", class_list(data, where), win_id, False, {})
    rows_out = []
    raw_rows = need(data, "rows", where)
    if not isinstance(raw_rows, list) or not raw_rows:
        die(f"rows manquant ({where})")
    flat_columns = []
    for row in raw_rows:
        reject_unknown(row, ROW_KEYS, where)
        row_id = str(need(row, "id", where))
        if row_id in ids:
            die(f"identifiant en double : {row_id}")
        ids.add(row_id)
        lines = need(row, "lines", row_id)
        if not isinstance(lines, int) or lines < 1:
            die(f"lines invalide ({row_id})")
        row_tag = str(need(row, "tag", row_id))
        if row_tag not in TAGS:
            die(f"balise inconnue « {row_tag} » ({row_id})")
        columns = []
        for col in need(row, "columns", row_id):
            reject_unknown(col, COL_KEYS, row_id)
            col_id = str(need(col, "id", row_id))
            if col_id in ids:
                die(f"identifiant en double : {col_id}")
            ids.add(col_id)
            tag = str(need(col, "tag", col_id))
            if tag not in TAGS or tag == "window":
                die(f"balise inconnue « {tag} » ({col_id})")
            fr = col.get("fr", 1)
            if not isinstance(fr, int) or fr < 1:
                die(f"fr invalide ({col_id})")
            item = {
                "id": col_id,
                "tag": tag,
                "fr": fr,
                "classes": class_list(col, col_id),
                "label": resolve_ref(str(col.get("label", "")), col_id),
                "file": str(col.get("file", "")),
                "name": str(col.get("name", "")),
                "save": str(col.get("save", "")),
                "run": "",
                "output": "",
                "readonly": True if "readonly" not in col else bool(col["readonly"]),
                "focus": parse_action(col.get("focus"), col_id),
                "enter": parse_action(col.get("enter"), col_id),
                "items": [],
            }
            if "\t" in item["label"] or "\n" in item["label"]:
                die(f"label illégal ({col_id})")
            if tag in ("nav", "select"):
                if "label" not in col or item["label"] == "":
                    die(f"label manquant ({col_id})")
                raw_items = need(col, "items", col_id)
                if not isinstance(raw_items, list):
                    die(f"liste invalide ({col_id})")
                item["items"] = [parse_item(raw, ids, col_id) for raw in raw_items]
            elif "items" in col:
                die(f"attribut items refusé ({col_id})")
            if tag == "input":
                item["file"] = str(need(col, "file", col_id))
                item["name"] = str(need(col, "name", col_id))
                item["save"] = str(need(col, "save", col_id))
                if item["save"] not in ("enter", "button"):
                    die(f"save invalide ({col_id})")
                if not item["enter"]["action"]:
                    item["enter"] = edit_self(col_id)
            if tag == "textarea":
                item["file"] = str(need(col, "file", col_id))
                if "readonly" not in col:
                    die(f"attribut manquant « readonly » ({col_id})")
                item["readonly"] = bool(col["readonly"])
                item["save"] = str(col.get("save", "button"))
                if item["save"] not in ("enter", "button"):
                    die(f"save invalide ({col_id})")
                if not item["readonly"] and not item["enter"]["action"]:
                    item["enter"] = edit_self(col_id)
            if tag == "button":
                if "label" not in col or item["label"] == "":
                    die(f"label manquant ({col_id})")
                if not item["enter"]["action"]:
                    die(f"action manquante ({col_id})")
                if item["enter"]["action"] not in ("save", "refresh"):
                    die(f"action de bouton inconnue ({col_id})")
            if tag == "command":
                if "label" not in col or item["label"] == "":
                    die(f"label manquant ({col_id})")
                if item["enter"]["action"] != "run":
                    die(f"commande manquante ({col_id})")
                item["run"] = item["enter"]["run"]
                item["output"] = item["enter"]["output"]
            if tag == "header" and ("label" not in col or item["label"] == ""):
                die(f"label manquant ({col_id})")
            item["focus"] = finish_set(item["focus"], "", col_id)
            item["enter"] = finish_set(item["enter"], "", col_id)
            columns.append(item)
            flat_columns.append(item)
        rows_out.append({
            "id": row_id,
            "tag": row_tag,
            "classes": class_list(row, row_id),
            "lines": lines,
            "columns": columns,
        })
    ensure_return(flat_columns, ids, where)
    for col in flat_columns:
        if col["file"] and not (ROOT / col["file"]).is_file():
            die(f"fichier introuvable : {col['file']}")
    hrefs = check_targets(flat_columns)
    win = {
        "id": win_id,
        "title": resolve_ref(str(data["label"]), where) if "label" in data else win_id,
        "file": path.name,
        "style": win_style,
        "rows": rows_out,
    }
    place_window(win, rules)
    windows.append(win)
    for href in hrefs:
        parse_window(path.parent / href, rules, windows, seen_paths)


def fitted_width(css_max, where):
    raw = (os.environ.get("COLUMNS") or "").strip()
    if not raw:
        return css_max
    if not raw.isdigit() or int(raw) < 1:
        die(f"COLUMNS invalide « {raw} » ({where})")
    return min(css_max, int(raw))


def place_window(win, rules):
    inherit = {"color": win["style"]["color"], "background": win["style"]["background"]}
    max_w = fitted_width(win["style"]["max_w"], win["id"])
    pad = win["style"]["pad"]
    inner = max_w - pad[1] - pad[3] - win["style"]["bw"] * 2
    y = 0
    prev_mb = None
    for row in win["rows"]:
        row_style = resolve(rules, row["tag"], row["classes"], row["id"], False, inherit)
        styled = []
        for col in row["columns"]:
            base = resolve(rules, col["tag"], col["classes"], col["id"], False, inherit)
            focus = resolve(rules, col["tag"], col["classes"], col["id"], True, inherit)
            styled.append((col, base, focus))
        gap = max([row_style["gap"]] + [st["gap"] for _c, st, _f in styled])
        frs = [col["fr"] for col, _b, _f in styled]
        widths = column_widths(inner, frs, gap, row["id"])
        heights = []
        placed = []
        x = pad[3] + win["style"]["bw"]
        for (col, base, focus), width in zip(styled, widths):
            tw = text_w(width, base)
            if tw < 1:
                die(f"zone trop étroite : {col['id']}")
            oh = outer_h(base, row["lines"])
            heights.append(oh)
            placed.append({
                "col": col,
                "base": base,
                "focus": focus,
                "x": x,
                "outer_w": width,
                "text_w": tw,
                "own_h": oh,
                "gap": gap,
            })
            x += width
            x += gap
        row_h = max(heights)
        mt = max(st["mar"][0] for _c, st, _f in styled)
        mb = max(st["mar"][2] for _c, st, _f in styled)
        if prev_mb is None:
            y += mt
        else:
            y += max(prev_mb, mt)
        row["y"] = y
        row["outer_h"] = row_h
        row["placed"] = placed
        y += row_h
        prev_mb = mb
        used = sum(widths) + gap * (len(widths) - 1)
        if used != inner:
            die(f"la rangée {row['id']} ne fait pas 100 % ({used} au lieu de {inner})")
    win["total_h"] = y + (prev_mb or 0)
    win["max_w"] = max_w
    win["inner"] = inner


def emit(windows, dest):
    lines = []
    for win in windows:
        lines.append("\x1f".join([
            "W", win["id"], str(win["max_w"]), win["title"], win["file"],
            win["style"]["background"], win["style"]["color"], str(win["total_h"]),
            str(win["inner"]),
        ]))
        order = 0
        for row in win["rows"]:
            cols = row["placed"]
            for index, zone in enumerate(cols):
                col = zone["col"]
                base = zone["base"]
                focus = zone["focus"]
                gap = zone["gap"] if index < len(cols) - 1 else 0
                content_h = row["outer_h"] - base["bw"] * 2
                lines.append("\x1f".join([
                    "Z", win["id"], col["id"], row["id"], str(index), col["tag"],
                    str(zone["x"]), str(row["y"]), str(zone["outer_w"]), str(row["outer_h"]),
                    str(zone["text_w"]), str(row["lines"]), str(content_h),
                    base["color"], base["background"], base["bc"], str(base["bw"]), str(base["bold"]),
                    focus["color"], focus["background"], focus["bc"], str(focus["bold"]),
                    str(base["pad"][0]), str(base["pad"][2]), str(base["pad"][3]), str(base["pad"][1]),
                    str(base["lh"]), str(gap),
                ]))
                props = {
                    "label": col["label"],
                    "file": col["file"],
                    "name": col["name"],
                    "save": col["save"],
                    "run": col["run"],
                    "output": col["output"],
                    "readonly": "1" if col["readonly"] else "0",
                }
                for prefix in ("focus", "enter"):
                    act = col[prefix]
                    if not act["action"]:
                        continue
                    for key in ("action", "target", "file", "name", "href", "run", "output", "value"):
                        if act[key]:
                            props[f"{prefix}.{key}"] = act[key]
                for key, val in props.items():
                    if val == "":
                        continue
                    lines.append("\x1f".join(["P", win["id"], col["id"], key, val]))
                if col["tag"] in ("nav", "select"):
                    for item_index, item in enumerate(col["items"]):
                        fields = [
                            "N", win["id"], col["id"], str(item_index), item["id"], item["label"],
                        ]
                        for act in (item["focus"], item["enter"]):
                            fields.extend([
                                act["action"], act["target"], act["file"], act["name"],
                                act["href"], act["run"], act["output"],
                            ])
                        fields.append(item["enter"]["value"])
                        lines.append("\x1f".join(fields))
                if col["tag"] in ("nav", "select", "input", "button", "command", "textarea"):
                    lines.append("\x1f".join(["F", win["id"], str(order), col["id"]]))
                    order += 1
    dest.write_text("\n".join(lines) + "\n", encoding="utf-8")


def build():
    css_path = ROOT / "theme.css"
    if not css_path.is_file():
        die("theme.css manquant")
    rules = parse_css(css_path.read_text(encoding="utf-8"))
    windows = []
    parse_window(ROOT / "windows" / "main.json", rules, windows, set())
    if not windows:
        die("aucune fenêtre")
    emit(windows, ROOT / "layout.tsv")
    return windows


def self_test(windows):
    def fr_widths(inner, frs, gap):
        return column_widths(inner, frs, gap, "essai")

    assert fr_widths(10, [1], 0) == [10]
    assert fr_widths(10, [1, 1, 2], 0) == [2, 2, 6]
    assert fr_widths(10, [1, 1, 2], 1) == [2, 2, 4]
    raw = (os.environ.get("COLUMNS") or "").strip()
    expected = 260
    if raw.isdigit() and int(raw) > 0:
        expected = min(260, int(raw))
    main = windows[0]
    assert main["max_w"] == expected
    header = main["rows"][0]
    assert len(header["placed"]) == 1
    assert header["placed"][0]["outer_w"] == main["inner"] == expected
    nav_row = main["rows"][1]
    widths = [z["outer_w"] for z in nav_row["placed"]]
    gaps = 1 * (len(widths) - 1)
    assert sum(widths) + gaps == expected, (widths, gaps)
    assert widths[-1] >= widths[0]
    zone = nav_row["placed"][0]
    assert zone["own_h"] == 6, zone["own_h"]
    assert nav_row["y"] == 5, nav_row["y"]
    labels = [z["col"]["label"] for z in nav_row["placed"]]
    assert labels == ["Rubriques", "Fiches", "Champs", "Menu"], labels
    champs = nav_row["placed"][2]["col"]
    assert champs["tag"] == "select"
    assert champs["items"][0]["focus"]["action"] == "load"
    assert champs["items"][0]["focus"]["name"] == "NOM"
    assert champs["items"][0]["enter"]["action"] == "edit"
    assert champs["items"][0]["enter"]["target"] == "valeur-hote"
    fiches = nav_row["placed"][1]["col"]
    assert fiches["items"][0]["label"] == "Passage"
    assert fiches["items"][0]["value"] == "127.0.0.1"
    assert fiches["items"][0]["enter"]["action"] == "set"
    assert fiches["items"][0]["enter"]["name"] == "HOTE"
    assert fiches["items"][1]["label"] == "Archive"
    jeton = fiches["items"][2]
    assert jeton["id"] == "fiche-jeton"
    assert jeton["label"] == "Jeton pro"
    assert jeton["value"] == "jeton-demo"
    assert jeton["focus"]["action"] == "load"
    assert jeton["focus"]["target"] == "valeur-hote"
    assert jeton["focus"]["file"] == "samples/passage.env"
    assert jeton["focus"]["name"] == "HOTE"
    assert jeton["enter"]["action"] == "set"
    assert jeton["enter"]["target"] == "valeur-hote"
    assert jeton["enter"]["file"] == "samples/passage.env"
    assert jeton["enter"]["name"] == "HOTE"
    assert jeton["enter"]["value"] == "jeton-demo"
    absolu = (ROOT / "samples" / "secrets.conf").resolve()
    assert absolu.is_absolute() and absolu.is_file()
    assert resolve_ref(f"{absolu}::HFU_PRO_JETON::label", "absolu") == "Jeton pro"
    assert resolve_ref(f"{absolu}::HFU_PRO_JETON::value", "absolu") == "jeton-demo"
    apercu = main["rows"][2]["placed"][1]["col"]
    assert apercu["id"] == "apercu" and apercu["readonly"] is False
    back = [it for it in zone["col"]["items"] if it["label"] == "Retour"]
    assert back and back[0]["enter"]["action"] == "return"
    assert back[0]["focus"]["action"] == ""
    assert any(w["file"] == "notes.json" for w in windows)
    print("self-test ok", labels)


def main():
    windows = build()
    if "--self-test" in sys.argv:
        self_test(windows)


if __name__ == "__main__":
    main()
