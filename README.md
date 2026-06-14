# YPlaza - Plateforme Immobilière

YPlaza est une plateforme web de gestion immobilière développée en Go, connectée à MySQL, avec un système d'enchères en temps réel, un algorithme de prévision des prix par localisation vérifiée (Nominatim / OpenStreetMap), une carte interactive Leaflet, et une administration centralisée.

📄 Un document de support technique détaillé (format .odt) accompagne ce README : `Support_Technique_YPlaza.odt` — il couvre le cahier des charges, les réalisations, les choix techniques et leur justification.

---

## Pourquoi Go ?

Go (Golang) a été choisi pour ce projet pour plusieurs raisons :

| Critère | Bénéfice |
|---|---|
| **Performance** | Compilation native, démarrage instantané, faible empreinte mémoire |
| **Simplicité** | Syntaxe claire, pas de frameworks lourds, déploiement en un binaire |
| **Concurrence** | Goroutines pour gérer les enchères, les sessions et les API externes en parallèle |
| **Standard library** | `net/http`, `text/template`, `database/sql` — pas de dépendances superflues |
| **Cross-platform** | Compilation pour Windows, Linux, macOS depuis une même base |
| **Personnel** | C'est l'un des principaux langages que je maîtrise, ce qui garantit un code maintenable |

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Frontend HTML/CSS/JS                   │
│  Templates Go (server-side rendering)                     │
│  Leaflet.js (cartes OpenStreetMap gratuites)              │
└──────────────────────┬──────────────────────────────────┘
                       │ HTTP
┌──────────────────────▼──────────────────────────────────┐
│              net/http (serveur natif)                     │
│  ├── Routes (/red_project/*)                              │
│  ├── Middleware (rate-limit, sessions)                    │
│  └── File server (/assets/*)                              │
└──────────────────────┬──────────────────────────────────┘
                       │ Appels
┌──────────────────────▼──────────────────────────────────┐
│               Controllers (handlers)                      │
│  ├── Web / Connect / Profil / Logement                    │
│  ├── Admin / Auction / Analytics                          │
└──────────────────────┬──────────────────────────────────┘
                       │ Requêtes
┌──────────────────────▼──────────────────────────────────┐
│              Database (couche d'accès)                    │
│  ├── user_db.go / property_db.go / auction_db.go          │
│  └── analytics_db.go / migrations                         │
└──────────────────────┬──────────────────────────────────┘
                       │ SQL + API externe
┌──────────────────────▼──────────────────────────────────┐
│           MySQL 8 (base "yplaza")                         │
│  Tables : users, properties, auctions, bids,              │
│           sales_history, user_liked_properties,            │
│           user_bought_properties                          │
├──────────────────────────────────────────────────────────┤
│           Services externes (gratuits, sans token)        │
│  ├── Nominatim API (OpenStreetMap) → géocodage           │
│  │   - Vérification code postal / ville / pays            │
│  │   - Coordonnées GPS pour la carte                      │
│  └── Leaflet.js + tuiles OpenStreetMap → carte            │
└──────────────────────────────────────────────────────────┘
```

## Stack technique

- **Backend** : Go 1.25, `net/http`, `golang.org/x/crypto/bcrypt`
- **Base de données** : MySQL 8 (via `go-sql-driver/mysql`)
- **Frontend** : HTML5, CSS3 (variables CSS, Flexbox/Grid), JavaScript vanilla
- **Templating** : `text/template` (server-side rendering) avec fonctions personnalisées (`add`, `mul`, `div`, `round`)
- **Cartographie** : Leaflet.js + tuiles OpenStreetMap (gratuit, sans clé API)
- **Géocodage** : Nominatim API (OpenStreetMap, gratuit, rate-limit 1 req/s)
- **Sécurité** : hash bcrypt, rate-limiting (5 échecs → verrouillage 15 min), confirmation email (token 2 minutes)
- **Déploiement** : Un seul binaire Go + MySQL

## Docker ?

Docker n'est pas utilisé actuellement, mais voici la réflexion :

| Pour Docker | Contre Docker |
|---|---|
| Uniformise l'environnement (Go + MySQL) | Surcharge pour un dev solo |
| Facilite le déploiement sur un VPS | Nécessite Docker Desktop (Windows) |
| Scaling horizontal avec plusieurs conteneurs | Le projet tient dans un binaire + MySQL |

**Conclusion** : Docker serait pertinent si l'application devait être déployée chez un hébergeur ou si plusieurs développeurs travaillaient dessus. Pour un projet étudiant en développement local, le lancement direct via `go run .` est plus simple et plus rapide. Un `docker-compose.yml` pourra être ajouté ultérieurement si besoin.

## Prérequis

- Go 1.25+
- MySQL 8 avec une base `yplaza`

## Installation et lancement

```bash
# 1. Configurer MySQL (variables d'env optionnelles, valeurs par défaut ci-dessous)
$env:DB_HOST="127.0.0.1"
$env:DB_PORT="3306"
$env:DB_USER="root"
$env:DB_PASS="motdepasse"
$env:DB_NAME="yplaza"

# 2. Lancer l'application
cd "ProjecT Dev"
go run . --kill-port-8080
```

Le serveur démarre sur **http://localhost:8080/red_project/home**

Les tables sont créées automatiquement au premier lancement (migration automatique).

### Flags disponibles

| Flag | Description |
|---|---|
| `--kill-port-80` | Tue le processus sur le port 80 au lancement |
| `--kill-port-8080` | Tue le processus sur le port 8080 au lancement |

### Arrêt du serveur

- **Ctrl+C** → arrêt gracieux avec `http.Server.Shutdown` (timeout 5s)
- Si le port est occupé, le programme refuse de démarrer avec un message explicite

## Comptes pré-configurés

| Rôle | Nom | Email | Mot de passe |
|---|---|---|---|
| Admin | alexandre | alexandre.petitfrere@ynov.com | Ynov_123 |
| Utilisateur | (à créer) | — | — |

## Routes principales

### Pages
| Route | Description |
|---|---|
| `/red_project/home` | Accueil avec liste des biens + filtres |
| `/red_project/logement/` | Détail d'un bien avec carte interactive |
| `/red_project/Profil` | Profil utilisateur |
| `/red_project/admin/users` | Administration (admin only) |
| `/red_project/auctions` | Liste des enchères en cours |
| `/red_project/analytics` | Analyses et prévisions avec estimation |

### Actions
| Route | Méthode | Description |
|---|---|---|
| `/red_project/login` | POST | Connexion (pseudo ou email) |
| `/red_project/register` | POST | Inscription avec email |
| `/red_project/confirm` | GET | Confirmation email (token) |
| `/red_project/logout` | POST | Déconnexion |
| `/red_project/logement/buy` | POST | Achat d'un bien |
| `/red_project/logement/post` | POST | Publier une annonce |
| `/red_project/logement/note` | POST | Noter un bien |
| `/red_project/auction/bid` | POST | Placer une enchère |
| `/red_project/auction/create` | POST | Créer une enchère |
| `/red_project/analytics/predict` | POST | Prédiction de prix (JSON) |
| `/red_project/analytics/geocode` | POST | Géocodage Nominatim (JSON) |
| `/red_project/analytics/report` | GET | Rapport ventes (JSON) |

## Fonctionnalités

### 🔐 Authentification sécurisée
- Hash bcrypt des mots de passe
- Confirmation par email (lien valable 2 minutes, SMTP configurable ou fallback console)
- Rate-limiting : 5 échecs → verrouillage 15 minutes
- Login accepte pseudo OU email

### 🏠 Gestion des biens
- **9 types de biens** : Maison, Appartement, Studio, Loft, Villa, Maison de ville, Penthouse, Local commercial, Terrain
- Liste complète avec filtres (type, prix min/max)
- Détail avec photos, prix, surface, nombre de pièces, localisation
- Affichage du **prix au m²** sur chaque bien
- **Carte interactive Leaflet.js** sur chaque fiche bien
- Achat / mise en vente

### 🔨 Système d'enchères
- Création d'enchères avec prix de départ, pas minimal et durée
- Affichage du prix actuel, du nombre d'offres et du temps restant
- Surenchère automatique
- Historique des offres avec horodatage
- Clôture automatique des enchères expirées

### 📊 Analyses et prévisions
- **Algorithme hybride** : régression linéaire + prix au m² + facteurs de localisation
- **Géocodage Nominatim** (OpenStreetMap) : vérification ville / code postal / pays
- **13 zones de prix** : Paris, Lyon, Marseille, Bordeaux, Toulouse, Lille, Nice, Nantes, Strasbourg, Montpellier, Rennes, Sud, Côte d'Azur
- **9 types de biens** avec multiplicateurs spécifiques
- **Prix au m²** calculé à partir de l'historique des ventes
- **Barre de confiance** et fourchette de prix
- **Carte interactive** du lieu estimé (Leaflet.js)
- Statistiques : ventes totales, prix moyen, prix/m² moyen, tendance du marché, région active
- Rapport de ventes exportable en JSON

### 🗺️ Cartographie (Leaflet.js + OpenStreetMap)
- **100% gratuit** — aucune clé API nécessaire
- Tuiles OpenStreetMap via CDN (`unpkg.com`)
- Marqueurs de localisation sur les fiches bien et les estimations
- Géocodage inverse via Nominatim (1 requête/seconde max)

### 🌍 Géocodage Nominatim
- API REST publique d'OpenStreetMap — **gratuite, sans token**
- Rate-limiting intégré (1 requête/seconde)
- Vérification : code postal → ville, ville → coordonnées GPS
- Endpoint dédié : `/red_project/analytics/geocode`

### 🛡️ Administration
- Gestion des utilisateurs (CRUD)
- Promotion / rétrogradation admin
- Statistiques globales

## Types de biens disponibles

| Type | Label | Multiplicateur prix |
|---|---|---|
| `house` | Maison | ×1.0 |
| `apartment` | Appartement | ×0.7 |
| `studio` | Studio | ×0.5 |
| `loft` | Loft | ×0.85 |
| `villa` | Villa | ×1.8 |
| `townhouse` | Maison de ville | ×1.1 |
| `penthouse` | Penthouse | ×1.5 |
| `commercial` | Local commercial | ×1.3 |
| `land` | Terrain | ×0.3 |

## Zones de prix (algorithme)

| Zone | Prix de base | €/m² | Coefficient |
|---|---|---|---|
| Paris | 500 000 € | 10 000 € | ×1.5 |
| Côte d'Azur | 450 000 € | 7 500 € | ×1.4 |
| Nice | 350 000 € | 6 500 € | ×1.3 |
| Sud | 350 000 € | 5 800 € | ×1.25 |
| Lyon | 300 000 € | 5 500 € | ×1.2 |
| Bordeaux | 280 000 € | 5 000 € | ×1.15 |
| Montpellier | 270 000 € | 4 500 € | ×1.15 |
| Toulouse | 260 000 € | 4 200 € | ×1.1 |
| Marseille | 250 000 € | 4 500 € | ×1.1 |
| Nantes | 250 000 € | 4 000 € | ×1.0 |
| Lille | 240 000 € | 4 000 € | ×1.0 |
| Strasbourg | 230 000 € | 3 800 € | ×1.0 |
| Rennes | 220 000 € | 3 800 € | ×1.0 |
| France (défaut) | 200 000 € | 3 500 € | ×1.0 |

## Algorithme de prédiction

1. **Géocodage** : la ville est vérifiée via Nominatim (OSM) → coordonnées GPS + code postal
2. **Requête SQL** : recherche des ventes similaires (même type, même localisation)
3. **Si ≥ 3 ventes** : calcul du prix au m² moyen + prix par pièce → moyenne pondérée
4. **Si < 3 ventes** : fallback avec élargissement géographique (toutes localisations)
5. **Si toujours < 3** : fallback théorique basé sur les facteurs de localisation
6. **Intervalle de confiance** : écart-type / prix moyen (borné entre 30% et 95%)
7. **Résultat** : prix estimé, fourchette, prix/m², nombre de données, confiance, carte

## Frontend — Design Méditerranéen

- **Palette** : Terracotta (#D4764A), Olive (#7A9E7E), Sable (#F5F0E8), Océan (#3D5A80)
- **Polices** : Playfair Display (titres) + Inter (corps)
- **Inspiration** : SeLoger.com, Airbnb
- **Responsive** : adaptation mobile/tablette
- **Composants** : cartes arrondies, ombres chaudes, badges, barres de progression, timeline d'offres

## Gestion du port 8080

- Vérification de disponibilité au démarrage
- Option `--kill-port-8080` pour libérer le port automatiquement
- Arrêt gracieux via `os.Signal` (Ctrl+C)
- Timeout de shutdown de 5 secondes
- Exécutable précédent renommé automatiquement

## Structure du projet

```
Support_Technique_YPlaza.odt    # Document de support technique (.odt)
ProjecT Dev/
├── main.go                    # Point d'entrée, flags, graceful shutdown
├── assets/
│   ├── css/                   # Styles (components, layout)
│   │   ├── home.main.css      # Point d'entrée CSS
│   │   ├── components/        # connect.css, proprety.css
│   │   └── layout/            # header.css, home.css, footer.css
│   ├── js/connect.js          # Modal + forms AJAX
│   └── img/                   # Images
├── controllers/               # Handlers HTTP
│   ├── Web.controllers.go     # Home, Init templates, reload
│   ├── Connect.controllers.go # Register, Login, SMTP
│   ├── DB.controllers.go      # CheckUser, RequireLogin, ConfirmEmail
│   ├── Profil.controllers.go  # Profil, Admin CRUD
│   ├── Logement.controllers.go # Properties, Auction handlers
│   └── Analytics.controllers.go # Stats, Predict, Geocode
├── database/                  # Accès MySQL + migrations
│   ├── database.go            # InitDB, migrate, bcrypt helpers
│   ├── user_db.go             # CRUD utilisateurs
│   ├── property_db.go         # CRUD propriétés + search
│   ├── auction_db.go          # Enchères + offres
│   └── analytics_db.go        # Ventes, stats, prédiction
├── models/                    # Structures de données
│   ├── struct.go              # User, Property, Auction, Prediction...
│   ├── structPage.go          # Home page model
│   └── const.go               # Constantes (couleurs console)
├── services/                  # Services externes
│   └── geocode.go             # Nominatim API client (OSM)
├── routes/routes.go           # Déclaration des routes
└── templates/                 # Templates HTML (12 fichiers)
    ├── header.html, footer.html
    ├── home.html, connect.html, error.html
    ├── logement.html, post_logement.html
    ├── profil.html, admin_users.html
    ├── auctions.html, auction_detail.html, create_auction.html
    └── analytics.html
```

## API Externe — Nominatim / OpenStreetMap

Le projet utilise **Nominatim**, l'API de géocodage gratuite d'OpenStreetMap.

| Caractéristique | Détail |
|---|---|
| **URL** | `https://nominatim.openstreetmap.org/search` |
| **Token** | Aucun requis (gratuit) |
| **Rate limit** | 1 requête/seconde (respecté par `services/geocode.go`) |
| **User-Agent** | `YPlaza/1.0 (immobilier)` |
| **Données** | Monde entier, villes, codes postaux, pays |
| **Usage** | Vérification adresse + coordonnées GPS pour la carte Leaflet |

### Endpoint YPlaza
`POST /red_project/analytics/geocode`
- `city` : nom de la ville
- `postal_code` : code postal
- `country` : pays (défaut: France)

→ Retourne `VerifiedLocation` (ville, code postal, pays, latitude, longitude)
