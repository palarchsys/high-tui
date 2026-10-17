# Grille JSON

L'écran est décrit par un fichier JSON par fenêtre. Les couleurs, le padding, le margin et la largeur maximale sont dans `theme.css`. Gum dessine les cadres. Python lit la grille une fois au lancement.

## Lancer

```bash
bash tui/demo.sh
```

`q` quitte. Les flèches déplacent le focus. Dans une liste, haut et bas changent de ligne. Dans une zone qui déborde, haut et bas font défiler. L'action `focus` part quand la ligne ou la zone reçoit le focus. L'action `enter` part avec Entrée. Échap revient à la fenêtre précédente.

La largeur affichée est `max-width`. Si `COLUMNS` est plus petit, la grille prend la largeur du terminal. Une ligne trop longue revient à la ligne. Si le texte dépasse la hauteur, une barre de défilement apparaît dans le cadre. Le cadre ne change pas de taille.

## Hauteur d'une rangée

`lines` est le nombre de lignes de texte à l'intérieur du cadre.

```
entre     = (lines - 1) * (line-height - 1)
intérieur = padding-haut + lines + entre + padding-bas
extérieur = border + intérieur + border
```

Le margin est en dehors du cadre. Entre deux rangées, l'écart est le plus grand des deux margins.

## Largeur

Une seule colonne prend toute la largeur utile. Sinon chaque colonne a un poids `fr` : toutes sauf la dernière reçoivent la part entière, et la dernière complète jusqu'à 100 %.

## Syntaxe

Chaque objet a un `id` et un `tag`. Le libellé visible est `label`. Une action est un objet, placé dans `focus` ou dans `enter` :

```json
{ "action": "load", "target": "valeur-hote", "file": "samples/passage.env", "name": "HOTE" }
```

L'action est `load`, `edit`, `set`, `window`, `return`, `save`, `refresh` ou `run`. On n'écrit que les champs utiles : `target`, `file`, `name`, `href`, `run`, `output`, `value`.

`focus` charge, par exemple, la variable de la ligne choisie. `enter` ouvre l'édition, écrit une valeur, change de fenêtre, enregistre, recharge ou lance une commande. La même action de focus ne repart pas tant que le focus ne change pas.

Un élément de liste a `label` (ce qui est affiché) et, pour une écriture directe, `value`. L'action `set` écrit cette valeur dans la variable `name` du fichier `file` dès qu'on presse Entrée.

Tout attribut `label` ou `value` égal à `chemin::NOM::label` ou `chemin::NOM::value` est résolu au lancement. Le chemin est absolu, ou relatif à la racine du produit. Le label est le texte de la ligne `# label NOM texte`. La value est le membre droit de `NOM=`. Un texte qui n'a pas cette forme n'est pas modifié. Une erreur arrête le lancement si le fichier, la variable ou le label manque, ou si le dernier mot n'est ni `label` ni `value`.

## Balises

`window`, `header`, `nav`, `select`, `input`, `textarea`, `button`, `command`, `status`.

`nav` et `select` sont des listes à hauteur fixe. `label` est le titre, sur la première ligne. Les lignes suivantes montrent les éléments. `›` marque la ligne choisie. Si la liste dépasse la hauteur, elle défile dans le cadre et une barre suit la position. Chaque élément a `id`, `label`, puis `focus` et `enter`. Sans `return`, « Retour » est ajouté à la première liste.

Un `input` nomme un fichier `NOM=valeur` et une variable. `save` vaut `enter` ou `button`. Le bouton Enregistrer écrit les zones modifiées de la fenêtre active. Recharger relit les fichiers de cette fenêtre.

Un `textarea` affiche un fichier. `readonly` vaut `true` ou `false`. S'il est éditable, Entrée ouvre l'éditeur sur tout le texte, même ce qui dépasse le cadre. Haut et bas font défiler la zone.

Une `command` a un `label` et une action `enter` de type `run`. `output` est l'`id` d'un textarea. Pendant l'exécution, la boucle s'arrête et le statut affiche « En cours : » suivi du label.

## Fichiers

- `theme.css` — seul endroit des couleurs et des mesures
- `windows/` — une fenêtre par fichier JSON
- `engine.py` — lecture du CSS et du JSON, écrit `layout.tsv`
- `theme.sh` — appels Gum, sans couleur écrite en dur
- `model.sh` — focus, fichiers, navigation, commandes
- `blocks.sh` — cache des blocs déjà dessinés
- `screen.sh` — assemblage et peinture, sans Gum
- `demo.sh` — boucle clavier

## Un seul cadre

```bash
DRAW_ONCE=1 bash tui/demo.sh
```

Affiche un écran et quitte. Cette voie n'appelle pas `gum input` ni `gum write`.
