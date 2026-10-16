# Grille JSON

L'écran est décrit par un fichier JSON par fenêtre. Les couleurs, le padding, le margin et la largeur maximale sont dans `theme.css`. Gum dessine les cadres. Python lit la grille une fois au lancement.

## Lancer

```bash
bash tui/demo.sh
```

`q` quitte. Les flèches déplacent le focus. Dans une liste, haut et bas changent de ligne. Entrée active l'élément. Échap revient à la fenêtre précédente. Si `COLUMNS` est inférieur à la largeur maximale, un cadre affiche « Terminal étroit », puis la grille.

Chaque zone garde toujours la même taille. Un texte trop long est coupé, le cadre ne bouge pas.

## Hauteur d'une rangée

`lines` est le nombre de lignes de texte à l'intérieur du cadre.

```
entre     = (lines - 1) * (line-height - 1)
intérieur = padding-haut + lines + entre + padding-bas
extérieur = border + intérieur + border
```

Le margin est en dehors du cadre. Entre deux rangées, l'écart est le plus grand des deux margins.

## Largeur

La fenêtre fait exactement `max-width` caractères. Une seule colonne prend toute la largeur utile. Sinon chaque colonne a un poids `fr` : toutes sauf la dernière reçoivent la part entière, et la dernière complète jusqu'à 100 %.

## Balises

`window`, `header`, `nav`, `input`, `textarea`, `button`, `command`, `status`. Chaque objet a un `id`. Une action vise cet `id`.

Un `nav` a un `title` et des `items`. L'action `load` remplit la zone `target`. L'action `window` ouvre le fichier `href`. L'action `return` revient à la fenêtre précédente. Sans item `return`, « Retour » est ajouté à la première liste.

Un `input` nomme un fichier `NOM=valeur` et une variable. `save` vaut `enter` ou `button`. Le bouton Enregistrer écrit les zones modifiées de la fenêtre active. Recharger relit les fichiers de cette fenêtre.

Un `textarea` affiche un fichier. `readonly` vaut `true` ou `false`.

Une `command` a un `title`, un `run`, et éventuellement `output`, l'`id` d'un textarea. Pendant l'exécution, la boucle s'arrête et le statut affiche « En cours : » suivi du titre.

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
