# Démo Gum

Explorateur de fiches en colonnes, thème sombre. Le titre du cadre est `high-tui`.

## Lancer

```bash
bash tui/demo.sh
```

`q` quitte. Les flèches changent de ligne ou de colonne. Entrée sur Rubriques, Fiches ou Champs place le focus sur Menu sans avancer de colonne. Ouvrir passe alors à la colonne suivante : Rubriques vers Fiches, Fiches vers Champs, Champs vers Menu. Modifier change la valeur choisie en mémoire pour la session, sans réécrire `data.tsv`. Retour revient aux champs. Échap seul, sur le menu, revient aussi aux champs.

Si `COLUMNS` est inférieur à 80, un cadre affiche « Terminal étroit », puis l'écran habituel.

## Fichiers

- `theme.sh` — palette, largeurs, bordures arrondies, aides Gum
- `data.tsv` — colonnes `rubrique`, `fiche`, `champ`, `valeur`
- `demo.sh` — écran et clavier

## Un seul cadre

```bash
DRAW_ONCE=1 bash tui/demo.sh
```

Affiche un écran et quitte. Cette voie n'appelle pas `gum input`.
