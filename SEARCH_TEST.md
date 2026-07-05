# Fonctionnalité de Recherche - Test

## Endpoint de Recherche

**Route:** `GET /search?query=<terme>`

**Accessibilité:** Public (pas d'authentification requise)

**Comportement:** Recherche dans tous les champs (title, director, genre, description) et retourne tous les films correspondant

## Exemples

### Test 1: Chercher "Matrix"
```
GET http://localhost:8080/search?query=Matrix
```

### Test 2: Chercher "Nolan" (directeur)
```
GET http://localhost:8080/search?query=Nolan
```

### Test 3: Chercher "Sci-Fi" (genre)
```
GET http://localhost:8080/search?query=Sci-Fi
```

### Test 4: Chercher "1999" (année)
```
GET http://localhost:8080/search?query=1999
```

### Test 5: Chercher "dream" (dans la description)
```
GET http://localhost:8080/search?query=dream
```

## Utilisation dans test-api.http

La requête est déjà configurée:

```
### 3. Search Movies (Public - pas de token requis)
GET {{baseUrl}}/search?query=Matrix
```

Remplacez "Matrix" par n'importe quel terme de recherche.

## Notes

- La recherche est **case-insensitive** (grâce à `LIKE`)
- Elle cherche dans **tous les champs**: titre, réalisateur, genre, description
- Elle retourne **tous les résultats correspondants**
- Le paramètre `query` est **obligatoire** (erreur 400 si manquant)
- Pas d'authentification requise (route publique)
