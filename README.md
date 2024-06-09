# go-market

Le projet permet au utilisateur de saisir des produits (id, titre, description, prix, quantité) des clients (id, nom, prenom, telephone, adresse, email)

Chaque client peut commander un seul produit avec une quantité (id, idclient, idproduit, quantite, prix, dateAchat) après chaque commande, envoyer par mail

Écrire un programme qui permet d'afficher un menu :

> 1. Ajouter un produit
> 2. Afficher l'ensemble des produits
> 3. Modifier un produit
> 4. Supprimer un produit
> 5. Exporter l'ensemble des produits sous forme csv
> 6. Ajouter un client
> 7. Afficher l'ensemble des clients
> 8. Modifier un client
> 9. Exporter l'ensemble des clients sous forme csv
> 10. Effectuer une commande (envoi mail + génération pdf)
> 11. Exporter l'ensemble des commandes
> 12. Quitter

## Mettre en route le projet

### Prérequis

- go 1.22
- docker

### Installation

```bash
git clone
cd go-market/mariadb
docker compose up -d
```

La base de données est prête à être utilisée au niveau de l'adresse `localhost:3306`

### Lancer le projet

```bash
go run .
```