# high-tui

Interface de terminal. Chaque fenêtre est un fichier HTML et un fichier CSS. La génération lit ces fichiers une fois. L'exécutable affiche ensuite la fenêtre courante avec Bubble Tea. Il ne relit pas le HTML.

## Générer, puis lancer

Les commandes se lancent dans le dossier `tui/` :

```sh
go generate ./...
go build -o high-tui ./cmd/high-tui
./high-tui
```

`go generate` écrit `internal/ui/generated.go`. Ce fichier est versionné. On ne le modifie pas à la main. `./high-tui` ne relance pas la génération.

`./high-tui --once` affiche une image et quitte.

La largeur affichée est le `max-width` de la fenêtre, limité par `COLUMNS` quand cette variable est un entier positif, puis par la largeur réelle du terminal. Sans `COLUMNS`, le `max-width` CSS est le plafond. `COLUMNS` invalide arrête le programme. Une zone plus étroite qu'un caractère aussi : `erreur : zone trop étroite`.

## Fichiers

| Fichier | Rôle |
| --- | --- |
| `windows/main.html` | Fenêtre de départ |
| `windows/<id>.html` | Une fenêtre |
| `windows/<id>.css` | La feuille de cette fenêtre |

Un `data-enter-href` est un chemin relatif au dossier `windows/`, par exemple `notes.html`.

Chaque fichier est un document HTML. Il commence par `<!DOCTYPE html>`, puis `html`, `head` et `body`. `head` contient `meta charset="utf-8"`, `title` et `link rel="stylesheet"` vers la feuille sœur. Pas de `script`, pas de style en ligne. Les commentaires `<!-- … -->` sont ignorés. Un identifiant est unique dans sa fenêtre. L'identifiant de fenêtre est celui de la `table` racine. Une balise inconnue, un attribut inconnu, un attribut présent mais vide, ou un identifiant en double arrête la génération. `readonly` peut être présent sans valeur.

On ouvre le fichier dans un navigateur pour voir la grille et les textes. Le navigateur ne lit pas le disque. Le texte visible est un placeholder. L'exécutable ignore ce texte dès qu'un `data-file` ou un `data-label` est posé, et il lit le fichier ou la référence.

## Grille

`body` contient une `table`. Une `table` contient des `tr`. Un `tr` contient des `td`. Un `td` contient une zone, ou une autre `table`. La zone occupe toute la cellule. La grille imbriquée aussi.

Une zone est `header`, `nav`, `input`, `textarea`, `button` ou `footer`.

`header` contient un `h1`. `nav` contient un `h2` puis un `ul` de `li`. Le texte du titre est le libellé. Le texte du `li` est le libellé affiché par le navigateur. `data-label`, s'il est présent, est le libellé de l'exécutable. Un `li` peut contenir un `a` : le navigateur suit `href`, l'exécutable suit `data-enter`.

`button` a `type="button"` et son texte est le libellé. `footer` est la ligne de statut. `input` et `textarea` portent le placeholder dans `value` ou dans le corps.

`header` et `footer` ne prennent pas le focus. Toutes les autres zones le prennent, y compris un `textarea` en lecture seule, pour pouvoir le faire défiler.

Si la fenêtre n'a aucune action `return`, la génération ajoute « Retour » sur la première `nav`. L'identifiant est celui de la liste plus `-retour`. Entrée quitte la fenêtre. Le focus ne fait rien.

## Attributs

Les attributs de l'application sont des `data-*`. Un attribut absent ne déclenche rien. Un attribut présent mais vide est une erreur, sauf `readonly`.

| Attribut | Absent |
| --- | --- |
| `id` | obligatoire sur la table, la zone et le `li` |
| `class` | aucune classe |
| `data-fr` | `1`, sur le `td` |
| `data-lines` | `1`, sur le `tr`. Hauteur de texte des zones de la rangée |
| `data-file` | vide, aucun fichier |
| `name`, `data-name` | vide, aucune variable |
| `data-save` | pas d'écriture automatique. Présent : `enter` ou `button` |
| `readonly` | faux. Présent, l'édition est refusée |
| `data-items` | les `li` du HTML. Présent : fichier ou `cmd:…`, les `li` ne sont qu'un aperçu |
| `data-next` | le focus ne remplit pas une autre liste |
| `data-focus` | vide, aucune action au focus |
| `data-enter` | vide, Entrée ne fait rien |
| `data-label`, `data-value` | le texte visible |
| `data-focus-target`, `data-focus-file`, `data-focus-name`, `data-focus-href`, `data-focus-run`, `data-focus-output`, `data-focus-value` | vides |
| `data-enter-target`, `data-enter-file`, `data-enter-name`, `data-enter-href`, `data-enter-run`, `data-enter-output`, `data-enter-value` | vides |

Les actions connues sont `load`, `edit`, `set`, `window`, `return`, `save`, `refresh` et `run`.

| Action | Exige |
| --- | --- |
| `load` | `target` et `file`. `name` vide charge tout le fichier dans la zone. Sinon la variable est lue, et la zone retient ce fichier et ce nom |
| `edit` | `target`, une zone `input` ou `textarea` qui n'est pas en lecture seule |
| `set` | `file`, `name` et `value`. Écrit tout de suite. `target`, s'il est présent, est un `input` ou un `textarea` |
| `window` | `href` |
| `return` | rien. Revient à la fenêtre précédente |
| `save` | rien. Écrit les zones modifiées de la fenêtre active |
| `refresh` | rien. Relit les fichiers et les références, et oublie les modifications non enregistrées |
| `run` | `run`. `output`, s'il est présent, est l'identifiant d'un `textarea` |

Une action inconnue, une cible absente ou un fichier absent arrête la génération. Les références `chemin::NOM::label` restent brutes dans le fichier généré. L'exécutable les résout au lancement et à chaque `refresh`.

Exemple :

```html
<li id="item-reseau" data-focus="load" data-focus-target="valeur-hote"
    data-focus-file="samples/passage.env" data-focus-name="HOTE"
    data-enter="edit" data-enter-target="valeur-hote">Réseau</li>
```

## Listes dynamiques

`data-items` sur une `nav` charge ses lignes au lancement. `data-next` nomme la `nav` qui reçoit la liste suivante. Au focus d'une ligne, cette nav est remplacée, sa sélection revient au début, et les nav suivantes de la chaîne sont vidées. Le focus reste sur la ligne qui a chargé.

La source est un tableau JSON, ou un objet JSON par ligne. Une commande s'écrit `cmd:` suivi de la commande, sans trou pour une valeur. Le fichier et la commande impriment la même forme.

| Champ | Rôle |
| --- | --- |
| `label` | Texte affiché. Obligatoire |
| `value`, `file`, `name` | Complètent l'action de la nav quand son fichier, son nom ou sa valeur est vide |
| `next` | Chemin, `cmd:…`, ou tableau de lignes |
| `focus`, `enter` | Action facultative, qui remplace celle de la nav. Champs : `kind`, `target`, `file`, `name`, `href`, `run`, `output`, `value` |

Un `href` dont le `kind` est `window` ouvre une fenêtre, résolu comme `data-enter-href` depuis le dossier `windows/`.

Un fichier est lu tout de suite. Une commande lancée par le focus passe par l'anneau d'attente. Échap l'annule. Un résultat en retard est ignoré. `refresh` relit `data-items` et rejoue le `next` de la ligne focalisée. Une source absente, un JSON illisible ou un libellé manquant s'affiche dans le statut et laisse la nav cible vide.

```json
{"label": "Réseau", "next": [{"label": "Passage", "file": "samples/passage.env", "name": "HOTE", "value": "127.0.0.1"}]}
```

## Références de variables

Dans `label` ou `value`, le texte `chemin::NOM::label` ou `chemin::NOM::value` est remplacé au lancement. `NOM` est le nom à gauche du `=`. Le libellé est la fin de la ligne `# label NOM …`. Un chemin relatif part du dossier `tui/`. Un chemin absolu est utilisé tel quel. Un autre texte reste littéral. Un troisième mot autre que `label` ou `value` est une erreur. Un fichier, une variable ou un libellé manquant arrête le lancement.

## CSS

Le fichier porte le même nom que le HTML. Les sélecteurs sont une balise, `.classe`, `#id` et `:focus`, séparés par des virgules. Pas de descendant, pas de `@media`, pas de `flex`, pas de `px`, pas de `em`, pas de script.

Les balises stylées sont `html`, `body`, `table`, `tr`, `td`, `h1`, `h2`, `ul`, `li`, `a`, `header`, `nav`, `input`, `textarea`, `button` et `footer`. Le style de la fenêtre est celui de `body`.

Une propriété inconnue arrête la génération en nommant l'identifiant. `display`, `width`, `height`, `box-sizing`, `table-layout`, `border-collapse`, `vertical-align`, `list-style` et `font` sont acceptés pour l'aperçu et ne changent pas le dessin du terminal.

| Propriété | Absent |
| --- | --- |
| `color` | hérite, sinon `#e6edf3` sur la fenêtre |
| `background`, `background-color` | hérite, sinon `#0d1117` sur la fenêtre |
| `border` | raccourci `0` ou `1`, `solid` ou `none`, et une couleur |
| `border-width` | `0`. Valeur admise : `0` ou `1` |
| `border-style` | `solid` ou `none` |
| `border-color` | la couleur du texte |
| `border-radius` | `0`. `1ch` arrondit le cadre |
| `padding`, `margin` | `0`. Une à quatre longueurs, dans l'ordre CSS |
| `max-width` | `260ch` sur `body` |
| `column-gap` | `0`. L'écart d'une rangée est le plus grand écart de ses zones et de la rangée |
| `line-height` | `1`, un entier, sans `ch` |
| `font-weight` | `normal`. L'autre valeur est `bold` |
| `--title-color` | la couleur du texte. Elle colore la première ligne source d'un `header`, d'un `nav` et du nom d'un `input`. Pendant l'édition, elle colore aussi le contour d'un `input` ou d'un `textarea` |

Une longueur est un nombre seul ou un nombre suivi de `ch`. `:focus` ne s'applique qu'à la zone qui a le focus. Une règle sans `:focus` s'applique aussi à la zone focalisée, donc `--title-color` posé sur la balise reste au focus.

Les marges verticales de deux rangées voisines se replient sur la plus grande. La marge haute de la première rangée et la marge basse de la dernière restent. À gauche, le contenu est rentré de la marge gauche de la fenêtre, de sa bordure, de son padding gauche et de la plus grande marge gauche des zones directement dans la rangée. À droite, le calcul est symétrique. Entre deux colonnes, l'écart est le plus grand entre le `column-gap` et les marges des zones qui se font face. Une grille dans une cellule n'ajoute pas de marge sur le bord : ses zones gardent les leurs à l'intérieur. Le cadre occupe toute la largeur retenue.

La hauteur extérieure d'une zone est `bordure + padding haut + lines + (lines − 1) × (line-height − 1) + padding bas + bordure`. Une seule colonne prend toute la largeur intérieure. Sinon les colonnes sauf la dernière reçoivent `floor(utile × fr / somme)`, et la dernière reçoit le reste. `utile` retire les écarts. `fr` est `data-fr`. La rangée prend la plus grande hauteur de ses cellules. Une zone plus courte est étirée. Une grille plus courte reçoit le surplus sur sa dernière rangée. La marge haute et la marge basse d'une grille imbriquée se replient avec les rangées voisines de la grille qui la contient.

Dans le navigateur, les marges CSS s'additionnent. L'aperçu montre la structure et les placeholders, pas l'écart du terminal au caractère près.

## Touches

| Touche | Effet |
| --- | --- |
| Gauche, Droite | Zone focalisable précédente ou suivante |
| Haut, Bas | Dans une liste, change la ligne et lance son action de focus. Dans un champ qui déborde, défile, puis change de zone au bord |
| Entrée | Hors édition, action `enter`. Dans un `textarea`, retour à la ligne. Dans un `input`, enregistre et sort |
| Ctrl+S | Dans l'éditeur, écrit le fichier de la zone et sort. Hors édition, ignoré |
| Ctrl+C | Dans l'éditeur, annule et restaure le texte. Ne quitte pas le programme |
| Ctrl+Q | Seule combinaison qui quitte. Une commande en cours est annulée d'abord |
| Échap | Dans l'éditeur, annule et restaure le texte. Sinon, fenêtre précédente, ou « Déjà à la première fenêtre ». Pendant une commande, annule |
| `q` | Ne quitte pas. Dans un éditeur, écrit la lettre |

L'action de focus part au lancement, au changement de zone ou de ligne, et à l'ouverture ou au retour d'une fenêtre. Elle ne repart pas à chaque image.

Entrée ouvre une zone éditable. Dans un `textarea`, Entrée insère ensuite un retour à la ligne. Maj+Entrée fait de même quand le terminal l'envoie. Ctrl+S écrit le fichier de la zone, affiche « Enregistré » et sort. Ctrl+C, comme Échap, abandonne le texte saisi, affiche « Annulé » et sort. Ctrl+C ne quitte jamais le programme. Pendant l'édition d'un `input` ou d'un `textarea`, le contour prend `--title-color`, la même couleur que les titres. Les flèches déplacent le curseur. Haut et Bas font défiler quand le curseur atteint le bord. Dans un `input`, Entrée et Ctrl+S enregistrent la valeur et sortent. Maj+Entrée n'y insère pas de ligne. Un collage conserve les retours à la ligne dans un `textarea` et les retire dans un `input`.

Le bouton Enregistrer écrit les zones encore marquées « Modifié ». `set` écrit tout de suite.

Une commande bloque les touches, affiche `⠋`, `⠙`, `⠹` ou `⠸` puis « En cours : » et le libellé. Échap annule. La sortie remplace le `textarea` nommé, même en lecture seule.

## Statuts

« Prêt », « Enregistré », « Rechargé », « Annulé », « Modifié », « Terminé » suivi du code, « Déjà à la première fenêtre ».
