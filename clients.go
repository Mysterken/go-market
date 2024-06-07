package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func displayClient(db *sql.DB, id int) bool {
	// 1. Récupérer les informations du client
	var name string
	var surname string
	var phone string
	var address string
	var email string

	err := db.QueryRow("SELECT name, surname, phone, address, email FROM clients WHERE id = ?", id).Scan(&name, &surname, &phone, &address, &email)

	if err != nil {
		fmt.Println("Client non trouvé")
		return false
	}

	// 2. Afficher les informations du client
	fmt.Printf("Nom: %s\n", name)
	fmt.Printf("Prénom: %s\n", surname)
	fmt.Printf("Téléphone: %s\n", phone)
	fmt.Printf("Adresse: %s\n", address)
	fmt.Printf("Email: %s\n", email)

	return true
}

func addClient(db *sql.DB) {
	reader := bufio.NewReader(os.Stdin)

	// 1. Demander à l'utilisateur de saisir les informations du client
	fmt.Println("Saisir le nom du client:")
	name, _ := reader.ReadString('\n')

	fmt.Println("Saisir le prénom du client:")
	surname, _ := reader.ReadString('\n')

	fmt.Println("Saisir le numéro de téléphone du client:")
	phone, _ := reader.ReadString('\n')

	fmt.Println("Saisir l'adresse du client:")
	address, _ := reader.ReadString('\n')

	fmt.Println("Saisir l'email du client:")
	email, _ := reader.ReadString('\n')

	// 2. Insérer les informations du client dans la table clients
	_, err := db.Exec("INSERT INTO clients (name, surname, phone, address, email) VALUES (?, ?, ?, ?, ?)", name, surname, phone, address, email)

	if err != nil {
		fmt.Println("Erreur lors de l'insertion du client")
		return
	}

	// 3. Afficher un message de confirmation
	fmt.Println("Client ajouté avec succès!")
}

func displayClients(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM clients")

	if err != nil {
		fmt.Println("Erreur lors de la récupération des clients")
		return
	}

	defer rows.Close()

	// 2. Afficher les clients
	fmt.Println("Liste des clients:")
	w := tabwriter.NewWriter(os.Stdout, 1, 2, 1, ' ', 0)
	fmt.Fprintln(w, "ID\tNom\tPrénom\tTéléphone\tAdresse\tEmail")
	for rows.Next() {
		var id int
		var name string
		var surname string
		var phone string
		var address string
		var email string
		rows.Scan(&id, &name, &surname, &phone, &address, &email)

		// trim spaces from name, surname and address
		name = strings.TrimSpace(name)
		surname = strings.TrimSpace(surname)
		address = strings.TrimSpace(address)

		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n", id, name, surname, phone, address, email)
	}

	w.Flush()
}

func modifyClient(db *sql.DB) {
	var clientChosen = false

	for !clientChosen {
		fmt.Println("Saisir l'identifiant du client à modifier:")
		var id int
		fmt.Scan(&id)

		if displayClient(db, id) {

			fmt.Println("Voulez-vous modifier ce client? (O/N)")
			var choice string
			fmt.Scan(&choice)

			if choice == "O" || choice == "o" {
				reader := bufio.NewReader(os.Stdin)

				fmt.Println("Saisir le nouveau nom du client:")
				name, _ := reader.ReadString('\n')

				fmt.Println("Saisir le nouveau prénom du client:")
				surname, _ := reader.ReadString('\n')

				fmt.Println("Saisir le nouveau numéro de téléphone du client:")
				phone, _ := reader.ReadString('\n')

				fmt.Println("Saisir la nouvelle adresse du client:")
				address, _ := reader.ReadString('\n')

				fmt.Println("Saisir le nouvel email du client:")
				email, _ := reader.ReadString('\n')

				_, err := db.Exec("UPDATE clients SET name = ?, surname = ?, phone = ?, address = ?, email = ? WHERE id = ?", name, surname, phone, address, email, id)

				if err != nil {
					fmt.Println("Erreur lors de la mise à jour du client")
					return
				}

				fmt.Println("Client modifié avec succès!")
				clientChosen = true
			} else {
				fmt.Println("Opération annulée")
				clientChosen = true
			}

		}
	}
}

func exportClients(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM clients")

	if err != nil {
		fmt.Println("Erreur lors de la récupération des clients")
		return
	}

	defer rows.Close()

	file, err := os.Create("clients.csv")

	if err != nil {
		fmt.Println("Erreur lors de la création du fichier")
		return
	}

	defer file.Close()

	writer := csv.NewWriter(file)

	writer.Write([]string{"ID", "Nom", "Prénom", "Téléphone", "Adresse", "Email"})

	for rows.Next() {
		var id int
		var name string
		var surname string
		var phone string
		var address string
		var email string
		rows.Scan(&id, &name, &surname, &phone, &address, &email)

		writer.Write([]string{fmt.Sprintf("%d", id), name, surname, phone, address, email})
	}

	writer.Flush()

	fmt.Println("Exportation des clients terminée")
}
