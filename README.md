# Movie Catalog API - API REST

## Description

Une API REST complète pour gérer un catalogue de films avec authentification JWT et gestion des utilisateurs. Développée en Go avec le framework Gin et une base de données SQLite.

## Fonctionnalités

### Endpoints CRUD
- ✅ Lister tous les films (`GET /movies`)
- ✅ Récupérer un film par ID (`GET /movies/:id`)
- ✅ Rechercher des films (`GET /search?query=...`)
- ✅ Créer un film (`POST /movies`)
- ✅ Mettre à jour un film (`PUT /movies/:id`)
- ✅ Supprimer un film (`DELETE /movies/:id`)
- ✅ Supprimer plusieurs films (`DELETE /movies`)
- ✅ Mettre à jour plusieurs films (`PUT /movies/bulk`)

### Authentification
- ✅ Inscription (`POST /user/register`)
- ✅ Connexion (`POST /user/login`)
- ✅ JWT Bearer Token
- ✅ Middleware d'authentification
- ✅ Gestion des admins

### Documentation
- ✅ Swagger/OpenAPI UI (`/swagger/index.html`)
- ✅ Documentation interactive

### Logging
- ✅ Fichier `logs/info.log` pour les opérations
- ✅ Fichier `logs/error.log` pour les erreurs

## Installation

### Prérequis
- Go 1.21 ou supérieur
- SQLite3

### Étapes d'installation

```bash
# Cloner le projet
cd biblio_go

# Installer les dépendances
go mod download

# Compiler
go build

# Ou lancer directement
go run .
```

## Configuration

Le serveur démarre sur `http://localhost:8080`

### Base de données
- Fichier: `database.db`
- Tables:
  - `users` (id, name, email, password, isAdmin)
  - `movies` (id, title, director, genre, release_year, description)

### Admin par défaut
- Email: `admin@mail.com`
- Mot de passe: `admin123`

## Utilisation

### 1. Inscription
```bash
POST /user/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

### 2. Connexion
```bash
POST /user/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

Réponse:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 3. Utiliser l'API (avec authentification)
```bash
GET /movies
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 4. Recherche
```bash
GET /search?query=Inception
Authorization: Bearer <token>
```

### 5. Créer un film
```bash
POST /movies
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Inception",
  "director": "Christopher Nolan",
  "genre": "Sci-Fi",
  "release_year": 2010,
  "description": "A skilled thief..."
}
```

### 6. Bulk Operations
```bash
# Mettre à jour plusieurs films
PUT /movies/bulk
Authorization: Bearer <token>
Content-Type: application/json

[
  {
    "id": 1,
    "title": "Updated Title",
    "director": "Director",
    "genre": "Genre",
    "release_year": 2020,
    "description": "Description"
  }
]

# Supprimer plusieurs films
DELETE /movies
Authorization: Bearer <token>
Content-Type: application/json

[1, 2, 3]
```

### 7. Profil utilisateur
```bash
# Mettre à jour son profil
PUT /user/me
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "New Name",
  "email": "newemail@example.com",
  "password": "newpassword123"
}

# Admin: Mettre à jour un utilisateur
PUT /user/:id
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "name": "Updated Name",
  "email": "user@example.com",
  "password": "newpassword"
}
```

### 8. Administration
```bash
# Réinitialiser la base de données (Admin uniquement)
DELETE /reset
Authorization: Bearer <admin-token>
```

## Documentation Swagger

Accédez à la documentation interactive Swagger à:
```
http://localhost:8080/swagger/index.html
```

## Structure du projet

```
biblio_go/
├── main.go              # Point d'entrée
├── database/
│   └── db.go            # Configuration BDD
├── models/
│   └── movie.go         # Modèle Movie
├── handlers/
│   └── movie_handlers.go   # Handlers des requêtes
├── services/
│   └── movie_services.go   # Logique métier
├── repositories/
│   └── movie_repository.go # Accès aux données
├── security/
│   ├── auth.go          # Authentification JWT
│   └── middleware.go    # Middleware d'authentification
├── routes/
│   └── routes.go        # Définition des routes
├── logger/
│   └── logger.go        # Système de logging
├── docs/                # Documentation Swagger
└── logs/                # Fichiers de log
    ├── info.log
    └── error.log
```

## Gestion des erreurs

L'API retourne les codes HTTP appropriés:
- `200 OK` - Succès
- `201 Created` - Ressource créée
- `400 Bad Request` - Erreur de validation
- `401 Unauthorized` - Authentification requise
- `403 Forbidden` - Accès refusé (admin)
- `404 Not Found` - Ressource non trouvée
- `500 Internal Server Error` - Erreur serveur

## Logging

Tous les événements sont enregistrés dans deux fichiers:
- **logs/info.log** - Opérations normales
- **logs/error.log** - Erreurs et problèmes

Format: `[TIMESTAMP] [LEVEL] Message`

## Déploiement

Pour déployer sur un hébergeur comme Render ou Alwaysdata:

1. Créer un compte sur le service
2. Connecter votre dépôt Git
3. Configurer les variables d'environnement
4. Le service se chargera de compiler et lancer l'application

Voir les fichiers `request.http` pour des exemples de requêtes complètes.

## Tests

Utilisez le fichier `request.http` pour tester tous les endpoints avec REST Client ou Thunder Client.

## Auteur

Projet réalisé pour l'école ESGI - Cours Go et API REST

## Licence

MIT