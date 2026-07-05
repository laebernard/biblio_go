# Guide de Test de l'API

## Prérequis
- VS Code avec l'extension **REST Client** (par Thunder Client ou similaire)
- Le serveur Go en cours d'exécution sur `http://localhost:8080`

## Étapes de Test

### 1. Démarrer le serveur
```bash
cd c:\Users\Loren\OneDrive\Bureau\ESGI\GoLang\Golang-Project\movie-go\biblio_go
go run .
```

Le serveur devrait afficher:
```
[GIN-debug] Listening and serving HTTP on :8080
```

### 2. Utiliser test-api.http pour les tests

Le fichier `test-api.http` contient toutes les requêtes pré-configurées. Voici comment l'utiliser:

#### Étape 1: Health Check
Cliquez sur "Send Request" pour tester la connexion:
```
GET http://localhost:8080/
```
Résultat attendu:
```json
{"message": "DB connected ✅"}
```

#### Étape 2: Login Admin
Envoyez la requête "Login Admin" pour obtenir un token:
```
POST http://localhost:8080/user/login
```
Réponse (exemple):
```json
{"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}
```

**IMPORTANT**: Copiez le token entier et remplacez la variable `@token = ` dans le fichier test-api.http

#### Étape 3-5: Tester les endpoints protégés
Une fois le token configuré, vous pouvez tester:
- Get All Movies
- Get Movie by ID
- Create Movie
- Update Movie
- Delete Movie
- Update My Profile
- Protected Route
- Reset Database

## Checklist de Fonctionnement

- [ ] Health Check retourne 200
- [ ] Login Admin retourne un token JWT
- [ ] Get All Movies retourne 200 (protégé)
- [ ] Create Movie retourne 201 (protégé)
- [ ] Update Movie retourne 200 (protégé)
- [ ] Delete Movie retourne 200 (protégé)
- [ ] Update My Profile retourne 200 (protégé)
- [ ] Protected Route retourne 200 (protégé)
- [ ] Reset Database retourne 200 (admin protégé)

## Logs

Les logs sont sauvegardés dans:
- `logs/info.log` - Opérations normales
- `logs/error.log` - Erreurs

## Dépannage

### Erreur: "Invalid token format"
- Assurez-vous que le token est copié correctement
- Le format doit être: `Authorization: Bearer <token>`

### Erreur: "Missing token"
- Le token n'a pas été fourni
- Vérifiez que `@token = <votre_token>` est bien défini

### Erreur: 404 Not Found
- La route n'existe pas
- Vérifiez que le serveur est lancé
- Vérifiez l'URL de la requête

### Erreur: "cannot fetch movies"
- La base de données a un problème
- Vérifiez que `database.db` existe
- Redémarrez le serveur
