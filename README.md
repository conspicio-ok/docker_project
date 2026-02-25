# Construire et demarrer l'application :

  
## 0) Démarrer :

Construire les images :

    docker compose build

Démarrer les conteneurs :

    docker compose up -d

Vérifier l’état :

    docker compose ps

Voir les logs :

    docker compose logs -f

Arrêter l’application :

    docker compose down


## 1) Structure du projet

    |
    ├── docker-compose.yml
    ├── .env
    ├── db.sql
    ├── backend/
    │   └── Dockerfile
    ├── Dockerfile_front
    ├── secret/
    │   ├── MYSQL_ROOT_PASSWORD.secret
    │   └── MYSQL_PASSWORD.secret


# 2) Architecture Docker

L’application est composée de trois services :

## Service mysql

- Image : mysql:9.6
- Conteneur : mysql_db
- Port exposé : 3306 (interne uniquement)
- Volume : mysql_data
- Réseau : backend_network (interne)

Rôle :

- Stockage des données
- Initialisation automatique via db.sql
- Données persistantes grâce au volume

---

## Service backend

- Build : ./backend
- Conteneur : go_api
- Port exposé : 8080
- Dépend de : mysql
- Réseaux : backend_network et frontend_network

Rôle :

- API développée en Go
- Communication avec MySQL
- Fournit les données au frontend

---

## Service front

- Image : mouchou/front-frontend
- Conteneur : front
- Port exposé : 80
- Dépend de : backend
- Réseau : frontend_network

Rôle :

- Interface utilisateur
- Accessible via http://localhost
- Communique avec l’API backend

---

# 3) Réseaux

## backend_network

- Type : bridge
- internal: true

Permet la communication sécurisée entre backend et MySQL.
MySQL n’est pas accessible depuis l’extérieur.

## frontend_network

- Type : bridge

Permet la communication entre frontend et backend.

---

# 4) Volume

Volume utilisé : mysql_data

Permet :

- La persistance des données
- La conservation des données même après suppression des conteneurs

Vérification :

    docker volume ls

---

# 5) Tests

## Test communication backend → MySQL

    docker compose logs mysql
    docker compose logs backend

Si aucune erreur de connexion n’apparaît, la communication fonctionne.

## Test communication frontend → backend

Navigateur :

    http://localhost

Ou en ligne de commande :

    curl http://localhost:8080

---

## Test persistance des données

1. Ajouter des jeux dans le panier
2. Arrêter les conteneurs :
   
       docker compose down

3. Redémarrer :
   
       docker compose up -d

Si les articles sont toujours dans le panier, la persistance fonctionne correctement.

!!! ATTENTION !!!
NE PAS FAIRE ``` docker compose down -v ``` si on veut la persisatnce des donnees