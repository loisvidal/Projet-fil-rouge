# YPlaza — Plateforme Immobilière

YPlaza est une plateforme web de gestion immobilière développée en Go avec MySQL. Elle propose la consultation et publication de biens, un système d'enchères en temps réel, un algorithme de prévision des prix avec géocodage Nominatim (OpenStreetMap), une carte interactive Leaflet, une administration centralisée, et un mode sombre/clair.

---

## Stack technique

| Couche | Technologie |
|---|---|
| **Backend** | Go (net/http, text/template, database/sql) |
| **Base de données** | MySQL 8 (go-sql-driver/mysql) |
| **Frontend** | HTML5, CSS3 (Flexbox/Grid, variables CSS), JavaScript vanilla |
| **Templating** | text/template avec fonctions perso (add, mul, div, round) |
| **Cartographie** | Leaflet.js + tuiles OpenStreetMap (gratuit, sans clé) |
| **Géocodage** | Nominatim API (gratuit, rate-limit 1 req/s, User-Agent YPlaza/1.0) |
| **Sécurité** | bcrypt, rate-limiting (5 échecs → blocage 15 min) |
| **Déploiement** | Binaire unique Go + MySQL (Docker multi-stage disponible) |

## Pourquoi Go ?

- **Performance** : compilation native, démarrage instantané, faible empreinte mémoire
- **Simplicité** : syntaxe claire, stdlib suffisante (net/http, text/template), un seul binaire
- **Concurrence** : goroutines pour enchères, sessions, appels API externes
- **Cross-platform** : compilation Windows/Linux/macOS depuis une même base

## Prérequis

- Go 1.25+
- MySQL 8 avec une base `yplaza` (ou Docker)

## Installation et lancement

```bash
# 1. Cloner le dépôt
git clone <url> && cd "ProjecT Dev"

# 2. Configurer MySQL (variables d'env facultatives)
$env:DB_HOST="127.0.0.1"
$env:DB_PORT="3306"
$env:DB_USER="root"
$env:DB_PASS="motdepasse"
$env:DB_NAME="yplaza"

# 3. Lancer l'application
go run . --kill-port-8080
```

Le serveur démarre sur **http://localhost:8080/red_project/home**

Les tables sont créées automatiquement au premier lancement.

### Avec Docker

```bash
docker compose up --build
```

### Flags

| Flag | Description |
|---|---|
| `--kill-port-8080` | Tue le processus sur le port 8080 au lancement |
| `--seed` | Réinitialise la base de données avec des données de démonstration |

### Arrêt

**Ctrl+C** → arrêt gracieux (timeout 5s). Si le port est occupé, le programme refuse de démarrer avec un message explicite.

## Comptes pré-configurés

| Rôle | Identifiant | Mot de passe |
|---|---|---|
| Admin | Alexandre | En_78270 |
| Utilisateur | toto | titi123 |

L'inscription est ouverte : après création du compte, l'utilisateur est automatiquement connecté (pas de confirmation email).

## Routes principales

### Pages

| Route | Description |
|---|---|
| `/red_project/home` | Accueil — liste des biens, filtres (type, prix, surface), dashboard |
| `/red_project/logement/{id}` | Détail d'un bien avec galerie, carte Leaflet, prix au m² |
| `/red_project/Profil` | Profil utilisateur |
| `/red_project/Profil/{id}` | Profil d'un autre utilisateur |
| `/red_project/admin/users` | Administration des utilisateurs (admin only) |
| `/red_project/auctions` | Liste des enchères en cours |
| `/red_project/auction/{id}` | Détail d'une enchère avec historique des offres |
| `/red_project/analytics` | Analyses et prédictions |
| `/red_project/contact` | Formulaire de contact |
| `/red_project/logement/post` | Publier une annonce (upload d'image) |
| `/red_project/auction/create` | Créer une enchère |

### Actions

| Route | Méthode | Description |
|---|---|---|
| `/red_project/login` | POST | Connexion (pseudo ou email) |
| `/red_project/register` | POST | Inscription (auto-confirm + auto-login) |
| `/red_project/logout` | POST | Déconnexion |
| `/red_project/logement/buy` | POST | Achat d'un bien |
| `/red_project/logement/post` | POST | Publier une annonce |
| `/red_project/Profil/delete` | POST | Supprimer son propre compte |
| `/red_project/auction/bid` | POST | Placer une enchère |
| `/red_project/auction/create` | POST | Créer une enchère |
| `/red_project/analytics/predict` | POST | Prédiction de prix (JSON) |
| `/red_project/analytics/geocode` | POST | Géocodage Nominatim (JSON) |
| `/red_project/admin/users/edit` | POST | Modifier un utilisateur (admin) |
| `/red_project/admin/users/delete` | POST | Supprimer un utilisateur (admin) |
| `/red_project/admin/users/toggle-admin` | POST | Promouvoir/rétrograder (admin) |
| `/red_project/home/filter` | POST | Filtrage des biens (type, prix, surface) |
| `/red_project/home/search` | POST | Recherche textuelle |

## Fonctionnalités

### 🔐 Authentification

- Hash bcrypt des mots de passe (≥ 6 caractères)
- Login accepte pseudo OU email
- Auto-confirm + auto-login après inscription
- Rate-limiting : 5 échecs → verrouillage 15 minutes
- Suppression de compte utilisateur (pas pour les comptes hardcodés ID ≤ 2)

### 🏠 Gestion des biens

- **9 types** : house, apartment, studio, loft, villa, townhouse, penthouse, commercial, land
- Publication d'annonce avec **upload d'image** (fichier sauvegardé dans assets/img/)
- Liste complète avec **filtres** : type de bien, prix min/max, **surface min/max**
- Détail avec galerie photos, prix, prix au m², surface, pièces, localisation
- **Carte interactive Leaflet** sur chaque fiche bien
- Achat / mise en vente

### 🔨 Système d'enchères

- Création avec prix de départ, pas minimal, durée (en heures)
- Affichage du prix actuel, nombre d'offres, temps restant (compte à rebours)
- Surenchère automatique
- Historique des offres avec horodatage
- Clôture automatique des enchères expirées

### 📊 Analyses et prévisions

- **Algorithme hybride** : régression linéaire + prix au m² + facteurs de localisation
- **Géocodage Nominatim** (OpenStreetMap) : vérification ville / code postal / pays
- **13 zones de prix** : Paris, Lyon, Marseille, Bordeaux, Toulouse, Lille, Nice, Nantes, Strasbourg, Montpellier, Rennes, Sud, Côte d'Azur
- **Barre de confiance** et fourchette de prix
- **Carte interactive** du lieu estimé
- Rapport de ventes exportable en JSON

### 🗺️ Cartographie

- 100% gratuit — aucune clé API (Leaflet.js + tuiles OpenStreetMap via unpkg.com)
- Marqueurs de localisation sur les fiches bien et les estimations
- Géocodage inverse via Nominatim (1 requête/s max)

### 🛡️ Administration

- Gestion des utilisateurs (CRUD)
- Promotion / rétrogradation admin
- Statistiques globales

## Frontend — Design

### Palette

| Variable | Light | Dark |
|---|---|---|
| `--gold` | #C5A059 | #D4AF37 |
| `--azure` | #007FFF | #007FFF |
| `--white` | #FFFFFF | rgb(30,30,30) |
| `--black` | rgb(33,33,33) | #FFFFFF |
| `--grey-light` | rgb(248,248,248) | rgb(10,10,10) |
| `--grey` | rgb(127,127,127) | rgb(180,180,180) |

### Dark/Light mode

- Basculé via un bouton dans le header (stockage localStorage, clé `theme`)
- `data-theme="dark"` / `data-theme="light"` sur `<html>`
- Fallback : `@media (prefers-color-scheme: dark)` si pas de préférence utilisateur

### Hero

- Image de fond `villa-1.jpeg` avec overlay sombre
- Texte dans un wrapper glassmorphism : `backdrop-filter: blur(14px)` + fond `rgba(245,245,245,0.75)` en light, adapté en dark

### Composants

- Search bar flottante (cartes, ombre portée)
- Catégories avec images de fond par type + overlay gradient
- Cartes de biens avec badge "À vendre" / "Vendu"
- Galerie d'images avec scroll snap
- Modale de connexion/inscription animée
- Pages auth : login.html, register.html (standalone, design cohérent)

### Responsive

Toutes les pages sont adaptées aux écrans mobiles et tablettes via des media queries à 900px, 768px, 600px et 480px :

- **Header** : passe en colonne, navigation compacte
- **Search bar** : colonne sur mobile, translateY supprimé
- **Grilles** : 2 colonnes à 768px, 1 colonne à 480px
- **Galeries** : hauteur réduite (220px), scroll horizontal
- **Formulaires** : champs et boutons full-width
- **Tableaux** : overflow-x auto avec scroll
- **Pages auth** : padding réduit, header empilé

## Types de biens

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

1. **Géocodage** : la ville est vérifiée via Nominatim → coordonnées GPS + code postal
2. **Requête SQL** : ventes similaires (même type, même localisation)
3. **Si ≥ 3 ventes** : prix au m² moyen + prix par pièce → moyenne pondérée
4. **Si < 3 ventes** : élargissement géographique (toutes localisations)
5. **Si toujours < 3** : fallback théorique basé sur les facteurs de localisation
6. **Intervalle de confiance** : écart-type / prix moyen (borné 30%–95%)

## API Externe — Nominatim

| Caractéristique | Détail |
|---|---|
| **URL** | `https://nominatim.openstreetmap.org/search` |
| **Token** | Aucun (gratuit) |
| **Rate limit** | 1 req/s (respecté par services/geocode.go) |
| **User-Agent** | `YPlaza/1.0` |
| **Usage** | Vérification adresse + coordonnées GPS pour carte Leaflet |

### Endpoint YPlaza

`POST /red_project/analytics/geocode`
- `city`, `postal_code`, `country` (défaut: France)
- Retourne `VerifiedLocation` (ville, code postal, pays, latitude, longitude)

## Gestion du port

- Vérification de disponibilité au démarrage
- Option `--kill-port-8080` pour libérer le port automatiquement
- Arrêt gracieux via `os.Signal` (Ctrl+C), timeout 5s
- Exécutable précédent renommé automatiquement

## Structure du projet

```
ProjecT Dev/
├── main.go                          # Point d'entrée, flags, graceful shutdown, MIME types
├── assets/
│   ├── css/
│   │   ├── home.main.css            # Point d'entrée (@import de 7 fichiers)
│   │   ├── style.css                # Pages auth standalone (login, register)
│   │   ├── animations.css           # Shimmer, gold pulse, modal animations
│   │   ├── components/
│   │   │   ├── connect.css          # Modale de connexion/inscription
│   │   │   └── proprety.css         # Cartes, détails, enchères, grilles (896 lignes)
│   │   └── layout/
│   │       ├── global.css           # Reset, variables CSS, container, boutons
│   │       ├── header.css           # Header principal responsive
│   │       ├── home.css             # Hero, search bar, catégories, about
│   │       ├── footer.css           # Footer responsive
│   │       ├── logement.css         # Galerie détail, bid form, bids list
│   │       └── profil.css           # Page profil responsive
│   ├── js/
│   │   ├── connect.js               # Modale + forms AJAX
│   │   └── theme.js                 # Dark/light mode toggle
│   └── img/                         # Images des biens + uploads
├── controllers/                     # Handlers HTTP
│   ├── Web.controllers.go           # Home, init templates, reload
│   ├── Connect.controllers.go       # Register, Login, SMTP
│   ├── DB.controllers.go            # CheckUser, requireLogin, confirmEmail
│   ├── Profil.controllers.go        # Profil, admin CRUD, delete own account
│   ├── Logement.controllers.go      # Properties, filter, auction handlers
│   └── Analytics.controllers.go     # Stats, predict, geocode
├── database/                        # Accès MySQL + migrations
│   ├── database.go                  # InitDB, migrate, bcrypt helpers
│   ├── user_db.go                   # CRUD utilisateurs
│   ├── property_db.go               # CRUD propriétés + search + images
│   ├── auction_db.go                # Enchères + offres
│   └── analytics_db.go              # Ventes, stats, prédiction
├── models/                          # Structures de données
│   ├── struct.go                    # User, Property, Auction, Prediction...
│   ├── structPage.go                # Home page model
│   └── const.go                     # Constantes
├── services/
│   └── geocode.go                   # Nominatim API client (OSM)
├── routes/routes.go                 # Déclaration des routes
├── templates/                       # 18 templates HTML
│   ├── header.html, footer.html, connect.html
│   ├── home.html, error.html
│   ├── logement.html, post_logement.html
│   ├── proprety.html (partial)
│   ├── profil.html, admin_users.html
│   ├── auctions.html, auction_detail.html, create_auction.html
│   ├── analytics.html, contact.html, legal.html
│   └── login.html, register.html
├── Dockerfile                       # Multi-stage build
└── docker-compose.yml               # Go + MySQL
```
