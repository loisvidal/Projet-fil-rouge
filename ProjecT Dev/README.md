# YPlaza - Projet Fil Rouge

YPlaza est une application web de petites annonces immobilieres developpee en Go. Le projet permet d'afficher une liste de biens, de gerer une connexion utilisateur simple et de preparer des fonctionnalites autour de l'achat, la vente, la consultation de logements et les profils.

## Technologies

- Go 1.23
- Serveur HTTP natif avec `net/http`
- Templates HTML avec `text/template`
- Donnees stockees en JSON
- HTML, CSS et JavaScript vanilla

## Prerequis

- Go installe sur la machine
- Un terminal place a la racine du projet

Verifier l'installation de Go :

```bash
go version
```

## Installation et lancement

Depuis la racine du projet :

```bash
go run .
```

Le serveur demarre sur :

```text
http://localhost:8080/red_project/home
```

Les fichiers statiques sont servis depuis le dossier `assets/` via l'URL `/assets/`.

## Structure du projet

```text
.
|-- assets/          # Fichiers statiques : CSS, JavaScript, images
|-- controllers/     # Handlers HTTP et logique applicative
|-- Data/            # Donnees JSON de l'application
|-- models/          # Structures Go partagees
|-- repository/      # Dossier prevu pour la couche d'acces aux donnees
|-- routes/          # Declaration des routes HTTP
|-- templates/       # Templates HTML
|-- go.mod           # Module Go
`-- main.go          # Point d'entree du serveur
```

## Fonctionnalites actuelles

- Affichage de la page d'accueil
- Chargement des biens immobiliers depuis `Data/dataProperty.json`
- Affichage des cartes de biens avec nom, description, prix, statut et images
- Inscription d'un utilisateur dans `Data/userList.json`
- Connexion d'un utilisateur existant
- Chargement des templates depuis le dossier `templates/`
- Service des assets CSS, JS et images

## Routes principales

| Methode | Route | Description |
| --- | --- | --- |
| GET | `/red_project/home` | Affiche la page d'accueil |
| POST | `/red_project/home/filter` | Route prevue pour filtrer les biens |
| POST | `/red_project/search` | Route prevue pour la recherche |
| POST | `/red_project/register` | Inscrit un nouvel utilisateur |
| POST | `/red_project/login` | Connecte un utilisateur |
| GET | `/red_project/logement/` | Route prevue pour consulter un logement |
| POST | `/red_project/logement/buy` | Route prevue pour acheter un logement |
| POST | `/red_project/logement/post` | Route prevue pour publier un logement |
| POST | `/red_project/logement/note` | Route prevue pour noter un logement |
| GET | `/red_project/Profil` | Route prevue pour afficher le profil |
| GET | `/red_project/Profil/Consult/` | Route prevue pour consulter un autre profil |
| POST | `/red_project/Profil/mail` | Route prevue pour contacter un profil |
| GET | `/red_project/error` | Affiche la page d'erreur |

## Donnees

Les donnees sont actuellement stockees dans des fichiers JSON :

- `Data/dataProperty.json` : liste des biens immobiliers
- `Data/userList.json` : liste des utilisateurs

Exemple de bien immobilier :

```json
{
  "NameProperty": "Studio etudiant Bordeaux",
  "DescProprety": "Studio meuble ideal pour etudiant, proche universite.",
  "IdProperty": 3,
  "PriceProperty": 145000,
  "IsSellProperty": true,
  "IsLiked": false,
  "ImgProperty": ["studio1.jpg"]
}
```

## Notes de developpement

Certaines routes sont declarees mais pas encore implementees. C'est notamment le cas des filtres, de la recherche, des pages logement, de l'achat, de la publication d'annonce, de la notation et des profils.

Les mots de passe sont stockes en clair dans le fichier JSON. Pour une version de production, il faudra ajouter un hachage des mots de passe, une vraie gestion de session et une validation plus stricte des formulaires.

## Commandes utiles

Lancer l'application :

```bash
go run .
```

Verifier que le projet compile :

```bash
go build .
```

Formater le code Go :

```bash
gofmt -w .
```
