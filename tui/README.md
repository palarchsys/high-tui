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

Un `href` est un chemin relatif au dossier `windows/`, par exemple `notes.html`.

Le document n'a pas de `html`, `head`, `body`, `script`, ni de style en ligne. Les commentaires `<!-- … -->` sont ignorés. Un identifiant est unique dans sa fenêtre. L'identifiant de fenêtre est unique dans le programme. Une balise inconnue, un attribut inconnu, un attribut présent mais vide, ou un identifiant en double arrête la génération.

## Balises

`window` contient des rangées. Une rangée est `header`, `main` ou `status`. Une rangée contient des zones. Une zone est `header`, `nav`, `select`, `input`, `textarea`, `button`, `command` ou `status`.

`nav` et `select` contiennent seulement des `item`. Les autres balises n'ont pas d'enfant.

`header` et `status` ne prennent pas le focus. Toutes les autres zones le prennent, y compris un `textarea` en lecture seule, pour pouvoir le faire défiler.

Si la fenêtre n'a aucune action `return`, la génération ajoute « Retour » dans cette fenêtre, sur la première `nav`, sinon sur la première `select`. L'identifiant est celui de la liste plus `-retour`. Entrée quitte la fenêtre. Le focus ne fait rien.

## Attributs

Toutes les balises acceptent les mêmes attributs. Un attribut absent ne déclenche rien. Un attribut présent mais vide est une erreur.

| Attribut | Absent |
| --- | --- |
| `id` | obligatoire |
| `class` | aucune classe |
| `label` | vide. Obligatoire sur `header`, `nav`, `select`, `button`, `command` et `item` |
| `value` | vide |
| `fr` | `1` |
| `lines` | `1`. Sur une rangée, c'est la hauteur de texte de chaque zone |
| `file` | vide, aucun fichier |
| `name` | vide, aucune variable |
| `href` | vide |
| `run` | vide |
| `output` | vide |
| `target` | vide |
| `save` | pas d'écriture automatique. Présent : `enter` ou `button` |
| `readonly` | `false`. L'édition ne démarre que si l'action est `edit`. `true` refuse cette action |
| `focus` | vide, aucune action au focus |
| `enter` | vide, Entrée ne fait rien |
| `focus-target`, `focus-file`, `focus-name`, `focus-href`, `focus-run`, `focus-output`, `focus-value` | vides |
| `enter-target`, `enter-file`, `enter-name`, `enter-href`, `enter-run`, `enter-output`, `enter-value` | vides |

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
<item id="item-reseau" label="Réseau"
      focus="load" focus-target="valeur-hote"
      focus-file="samples/passage.env" focus-name="HOTE"
      enter="edit" enter-target="valeur-hote">
</item>
```

## Références de variables

Dans `label` ou `value`, le texte `chemin::NOM::label` ou `chemin::NOM::value` est remplacé au lancement. `NOM` est le nom à gauche du `=`. Le libellé est la fin de la ligne `# label NOM …`. Un chemin relatif part du dossier `tui/`. Un chemin absolu est utilisé tel quel. Un autre texte reste littéral. Un troisième mot autre que `label` ou `value` est une erreur. Un fichier, une variable ou un libellé manquant arrête le lancement.

## CSS

Le fichier porte le même nom que le HTML. Les sélecteurs sont une balise, `.classe`, `#id` et `:focus`, séparés par des virgules. Pas de descendant, pas de `@media`, pas de `flex`, pas de `px`, pas de `em`, pas de script.

`item` n'est pas une balise CSS. Les balises stylées sont `window`, `header`, `main`, `nav`, `select`, `input`, `textarea`, `button`, `command` et `status`.

Une propriété inconnue arrête la génération en nommant l'identifiant.

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
| `max-width` | `260ch` sur la fenêtre |
| `column-gap` | `0`. L'écart d'une rangée est le plus grand écart de ses zones et de la rangée |
| `line-height` | `1`, un entier, sans `ch` |
| `font-weight` | `normal`. L'autre valeur est `bold` |
| `title-color` | la couleur du texte. Elle colore la première ligne source d'un `header`, d'un `nav`, d'un `select` et du nom d'un `input`. Pendant l'édition, elle colore aussi le contour d'un `input` ou d'un `textarea` |

Une longueur est un nombre seul ou un nombre suivi de `ch`. `:focus` ne s'applique qu'à la zone qui a le focus. Une règle sans `:focus` s'applique aussi à la zone focalisée, donc `title-color` posé sur la balise reste au focus.

Les marges verticales de deux rangées voisines se replient sur la plus grande. La marge haute de la première rangée et la marge basse de la dernière restent. À gauche, le contenu est rentré de la marge gauche de la fenêtre, de sa bordure, de son padding gauche et de la plus grande marge gauche des zones de la rangée. À droite, le calcul est symétrique. Entre deux colonnes, l'écart est le plus grand entre le `column-gap` et les marges qui se font face. Le cadre occupe toute la largeur retenue.

La hauteur extérieure d'une zone est `bordure + padding haut + lines + (lines − 1) × (line-height − 1) + padding bas + bordure`. Une seule colonne prend toute la largeur intérieure. Sinon les colonnes sauf la dernière reçoivent `floor(utile × fr / somme)`, et la dernière reçoit le reste. `utile` retire les `column-gap`.

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

Entrée ouvre une zone éditable. Dans un `textarea`, Entrée insère ensuite un retour à la ligne. Maj+Entrée fait de même quand le terminal l'envoie. Ctrl+S écrit le fichier de la zone, affiche « Enregistré » et sort. Ctrl+C, comme Échap, abandonne le texte saisi, affiche « Annulé » et sort. Ctrl+C ne quitte jamais le programme. Pendant l'édition d'un `input` ou d'un `textarea`, le contour prend `title-color`, la même couleur que les titres. Les flèches déplacent le curseur. Haut et Bas font défiler quand le curseur atteint le bord. Dans un `input`, Entrée et Ctrl+S enregistrent la valeur et sortent. Maj+Entrée n'y insère pas de ligne. Un collage conserve les retours à la ligne dans un `textarea` et les retire dans un `input`.

Le bouton Enregistrer écrit les zones encore marquées « Modifié ». `set` écrit tout de suite.

Une commande bloque les touches, affiche `⠋`, `⠙`, `⠹` ou `⠸` puis « En cours : » et le libellé. Échap annule. La sortie remplace le `textarea` nommé, même en lecture seule.

## Statuts

« Prêt », « Enregistré », « Rechargé », « Annulé », « Modifié », « Terminé » suivi du code, « Déjà à la première fenêtre ».
