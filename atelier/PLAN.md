# Atelier high-tui — plan

Éditeur visuel à côté du générateur. `go generate` ne change pas. L'atelier produit les mêmes `windows/<id>.html` et la feuille sœur. Le programme lit ces fichiers une fois, puis n'y revient plus.

Le document interne de l'atelier est l'arbre déjà décrit par `model.Window` : grilles, zones, actions, lignes. Le HTML et le CSS sont une projection. Un document illégal est refusé avec les diagnostics du générateur (`erreur : …`).

## Ce qui ne bouge pas

Une fenêtre est une `table`. Un `tr` porte `data-lines`. Un `td` porte `data-fr` et contient une zone ou une grille. Zones : `header`, `nav`, `input`, `textarea`, `button`, `footer`. Pas de script, pas de style en ligne, pas de balise inconnue.

Actions : `load`, `edit`, `set`, `window`, `return`, `save`, `refresh`, `run`.

CSS admis : couleurs, `border` 0 ou 1, `border-radius`, `padding`, `margin`, `max-width`, `column-gap`, `line-height`, `font-weight`, `--title-color`. Unités `ch`. Pas de descendant, pas de `px`, pas de `flex`.

## Deux parties

Bureau : affichage à gauche, données à droite, prévisualisation en bande sous les deux. Écran étroit : onglets Affichage, Données, Aperçu. La souris prime : dépôt, échange, poignées `fr` et `data-lines`.

### Affichage

- Palette des six zones. Cliquer une zone arme le dépôt ; cliquer une rangée ajoute une colonne ; la bande du bas ajoute une rangée.
- Glisser une zone sur une autre les échange.
- Poignée entre deux colonnes pour `data-fr`. Bord de rangée pour `data-lines`.
- Imbriquer une cellule dans une grille. La légende de la grille reçoit les nouvelles rangées.
- Inspecteur : identifiant, textes, `readonly`, `data-save`, actions, style de zone et style de fenêtre (`body`, dont `max-width`).

### Données

Cinq origines : en dur, fichier, commande `cmd:`, ligne active, référence `chemin::NOM::label|value`.

| Mode | Split | Branchement |
| --- | --- | --- |
| Texte en dur | non | libellé ou placeholder dans le HTML |
| Liste en dur | non | `<li>`. Pas de `next` : le HTML ne porte pas ce champ |
| Fichier brut | non | `data-file` sans nom, sur un champ |
| Variable | env label/value | `data-file` + `name` |
| Référence | env | le texte devient `chemin::NOM::label` ou `::value` |
| Liste fichier | json natif | `data-items` |
| Liste env, lignes ou CSV | label/value | projetée en JSON, `data-items` pointe ce JSON, le schéma reste sur la source |
| Liste de scalaires | non | une ligne de texte = un label, même projection |
| Enregistrement | env | chaque nom se branche comme une variable |
| Commande liste | json imprimé | `cmd:…`, fixture collée pour l'aperçu |
| Commande texte | non | action `run` vers un `textarea` |
| Écriture | | `data-save` `enter` ou `button`, ou `set` |

Un fil relie un port à un slot. `brut` et `value` vont à un champ. `label` va à un libellé. `lignes` va à une nav. Un slot n'a qu'un fil. Le dépôt incompatible est refusé.

Filtres de lecture, pas des modes : ignorer les lignes vides, trier par label.

### next-data

Méthode déjà dans l'exécutable. `data-next` sur une nav nomme une autre nav. La cible doit exister, être une nav, et la chaîne ne boucle pas (`cycle de listes`).

Le `next` d'une ligne a une seule forme :

| Forme | Effet |
| --- | --- |
| absent | la cible est vidée |
| fichier | lu tout de suite, même JSON que `data-items` |
| `cmd:` | anneau d'attente ; Échap annule ; un résultat en retard est ignoré |
| tableau | installé tout de suite |

Au focus d'une nav qui a `data-next`, avant l'action `focus` :

1. La commande de liste en cours est annulée.
2. La cible est vidée, sélection à zéro, puis toute la chaîne derrière elle. Les `<li>` d'aperçu disparaissent là.
3. La ligne focalisée fournit le contenu. Le tableau gagne sur un chemin.
4. Les lignes sont clonées, y compris leurs `next`. Sans identifiant : `nav-1`, `nav-2`. Les références sont résolues sur le libellé, la valeur et les valeurs d'action.
5. La sélection de la cible repart à zéro. Plus loin, la chaîne reste vide tant que la cible n'a pas le focus.

`data-items` est l'autre chargement : la nav elle-même, au lancement et au `refresh`. Une nav peut avoir les deux. Une ligne `return` du HTML n'est conservée qu'au chargement `data-items`, et seulement si le fichier n'en apporte pas. next-data ne la réinjecte pas.

L'action de la ligne, sinon celle de la nav, part ensuite. `file`, `name` et `value` vides sont complétés par la ligne.

Dans l'atelier, le choix `data-next` refuse tout de suite une mauvaise cible. Le `next` s'édite sur les lignes du JSON. L'aperçu rejoue la chaîne dans l'ordre du document, selon la ligne choisie : c'est l'état une fois la chaîne parcourue. Une commande sans fixture affiche `⠋ En cours` et ne lance rien.

## Prévisualisation

Grille de caractères. Les marges verticales se replient. `fr`, `data-lines`, `column-gap` et la marge de `body` suivent les règles du terminal, pas celles du navigateur. Largeur = `max-width`, plafonnée par les colonnes choisies. Cliquer une zone la focalise. Cliquer une ligne change la sélection et rejoue next-data. Le premier `footer` affiche le statut.

## Export

Valider, puis télécharger `windows/<id>.html` et `.css`. `generated.go` n'est pas écrit ici. Importer un document du même dialecte. Le projet tient dans le navigateur. Revenir aux fenêtres `main` et `notes` restaure les exemples.

## Technique

React, document Zustand, annulation locale. Pas d'éditeur HTML générique. Glisser au pointeur sur la grille. Pas de compte, pas de base. Palette du terminal : fond `#0d1117`, surface `#161b22`, texte `#e6edf3`, titre `#ffa657`, filet `#30363d`.

## Hors de ce jet

Modifier le générateur. Exécuter un `cmd:` pour de vrai. CSS libre. Jointure de fichiers. Modèle de texte `Hôte : {value}`.
