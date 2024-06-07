package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

func main() {

	db := connectToDB()

	for {
		showMenu()

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			addProduct(db)
		case 2:
			displayProducts(db)
		case 3:
			modifyProduct(db)
		case 4:
			deleteProduct(db)
		case 5:
			exportProducts(db)
		case 6:
			addClient(db)
		case 7:
			displayClients(db)
		case 8:
			modifyClient(db)
		case 9:
			exportClients(db)
		case 10:
			makeOrder(db)
		case 11:
			exportOrders(db)
		case 12:
			fmt.Println("Au revoir!")
			db.Close()
			os.Exit(0)
		default:
			fmt.Println("Choix invalide! Veuillez réessayer.")
		}
	}
}

func connectToDB() *sql.DB {
	// connect to mariadb at localhost:3306
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/go_market")

	if err != nil {
		fmt.Println("Error connecting to database: ", err)
		os.Exit(1)
	}

	return db
}

func showMenu() {
	fmt.Println(`
	1. Ajouter un produit
	2. Afficher l'ensemble des produits
	3. Modifier un produit
	4. Supprimer un produit
	5. Exporter l'ensemble des produits sous forme csv
	6. Ajouter un client
	7. Afficher l'ensemble des clients
	8. Modifier un client
	9. Exporter l'ensemble des clients sous forme csv
	10. Effectuer une commande (envoi mail + génération pdf)
	11. Exporter l'ensemble des commandes
	12. Quitter
	Choisir une option:`)
}

func makeOrder(db *sql.DB) {
	// 1. Demander à l'utilisateur de saisir les informations de la commande

	var clientChosen bool = false
	var clientId int

	for !clientChosen {
		fmt.Println("Saisir l'identifiant du client:")

		fmt.Scan(&clientId)

		if displayClient(db, clientId) {
			// ask if the user wants to modify the client
			fmt.Println("Voulez-vous choisir ce client? (O/N)")
			var choice string
			fmt.Scan(&choice)

			if choice == "O" || choice == "o" {
				clientChosen = true
			} else {
				fmt.Println("Opération annulée")
				return
			}
		}
	}

	var productChosen bool = false
	var productId int

	for !productChosen {
		fmt.Println("Saisir l'identifiant du produit:")

		fmt.Scan(&productId)

		if displayProduct(db, productId) {
			// ask if the user wants to modify the client
			fmt.Println("Voulez-vous choisir ce produit? (O/N)")
			var choice string
			fmt.Scan(&choice)

			if choice == "O" || choice == "o" {
				productChosen = true
			} else {
				fmt.Println("Opération annulée")
				return
			}
		}

	}

	fmt.Println("Saisir la quantité du produit:")
	var quantity int
	fmt.Scan(&quantity)

	// 2. Insérer les informations de la commande dans la table orders
	_, err := db.Exec("INSERT INTO orders (client_id, product_id, quantity, price, purchase_date) VALUES (?, ?, ?, (SELECT price FROM products WHERE id = ?), NOW())", clientId, productId, quantity, productId)

	if err != nil {
		fmt.Println("Erreur lors de l'insertion de la commande")
		return
	}

	// 3. Afficher un message de confirmation
	fmt.Println("Commande effectuée avec succès!")

	// 4. Envoyer un mail au client

	// 5. Générer un PDF de la commande
}

func exportOrders(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM orders")

	if err != nil {
		fmt.Println("Erreur lors de la récupération des commandes")
		return
	}

	defer rows.Close()

	file, err := os.Create("orders.csv")

	if err != nil {
		fmt.Println("Erreur lors de la création du fichier")
		return
	}

	defer file.Close()

	writer := csv.NewWriter(file)

	writer.Write([]string{"ID", "Client ID", "Product ID", "Quantité", "Prix", "Date d'achat"})

	for rows.Next() {
		var id int
		var client_id int
		var product_id int
		var quantity int
		var price float64
		var purchase_date string
		rows.Scan(&id, &client_id, &product_id, &quantity, &price, &purchase_date)

		writer.Write([]string{fmt.Sprintf("%d", id), fmt.Sprintf("%d", client_id), fmt.Sprintf("%d", product_id), fmt.Sprintf("%d", quantity), fmt.Sprintf("%.2f", price), purchase_date})
	}

	writer.Flush()

	fmt.Println("Exportation des commandes terminée")
}
