package main

import (
    "encoding/json"
    "fmt"
    bolt "github.com/coreos/bbolt"
    "log"
    "strconv"
)

type Cliente struct {
    NroCliente int
    Nombre     string
    Apellido   string
    Domicilio  string
    Telefono   string
}

type Tarjeta struct {
    NroTarjeta   string
    NroCliente   int
    ValidaDesde  string
    ValidaHasta  string
    CodSeguridad int
    LimiteCompra float64
    Estado       string
}

type Comercio struct {
    NroComercio  int
    Nombre       string
    Domicilio    string
    CodigoPostal int
    Telefono     string
}

type Compra struct {
    NroOperacion int
    NroTarjeta   string
    NroComercio  int
    Fecha        string
    Monto        float64
    Pagado       bool
}

const separador string = "********************************************************************************\n"

var db *bolt.DB
var err error
var opcion int = 0
var clientes []Cliente
var tarjetas []Tarjeta
var comercios []Comercio
var compras []Compra

func main() {
    fillArrays()

    for opcion != 3 {
		fmt.Print("\n", separador)
		
        fmt.Print("\nMENU\n\n")
        fmt.Print("1 - Crear base de datos Tarjetas\n")
        fmt.Print("2 - Insertar datos a tablas\n") 
        fmt.Print("3-  Terminar y salir\n")

        fmt.Print("\nElija una opcion: ")
        fmt.Scanf("%d", &opcion)
        fmt.Print("Opcion seleccionada: ", opcion)
        fmt.Print("\n")
		
		fmt.Print("\n", separador)
		
        if opcion == 1 {
            db, err = bolt.Open("tarjetas.db", 0600, nil)
            if err != nil {
                log.Fatal(err)
            }
            defer db.Close()
            fmt.Print("\nSE CREO LA BASE DE DATOS\n\n", separador)
        }

        if opcion == 2 {
            insertData()
        }

        if opcion == 3 {
            fmt.Print("\nSALIR\n")
        }
    }
}

func fillArrays() {
	//Agregacion de 3 elementos por cada entidad
	clientes = append(clientes, Cliente{1, "Juan", "Riquelme", "Bransen 123", "541130888931"})
	clientes = append(clientes, Cliente{2, "Javier", "Saviola", "Nuñez 234", "541130888932"})
	clientes = append(clientes, Cliente{3, "Juan", "Cavaliere", "Ricchieri 950", "541130888933"})
	tarjetas = append(tarjetas, Tarjeta{"4540730040713558", 1, "202001", "202401", 356, 15000.0, "vigente"})
	tarjetas = append(tarjetas, Tarjeta{"4540730040713559", 2, "201901", "202301", 357, 10000.0, "vigente"})
	tarjetas = append(tarjetas, Tarjeta{"4540730040713560", 3, "201803", "202203", 358, 12000.0, "vigente"})
	comercios = append(comercios, Comercio{1, "Supermercado Dia", "Sdor Moron 200", 1661, "46660011"})
	comercios = append(comercios, Comercio{2, "Supermercado Disco", "Sdor Moron 1600", 1661, "46660022"})
	comercios = append(comercios, Comercio{3, "Supermercado Coto", "Munzon 6200", 1662, "46660033"})
	compras = append(compras, Compra{1, "4540730040713558", 356, "2021-06-14 18:39:45", 1000.0, false})
	compras = append(compras, Compra{2, "4540730040713559", 357, "2021-06-14 18:49:34", 5000.0, false})
	compras = append(compras, Compra{3, "4540730040713560", 358, "2021-06-14 18:55:26", 6000.0, false})
}

func insertData() {
    //insert a tabla cliente
    fmt.Print("\nTABLA CLIENTE\n\n")
    for _, cliente := range clientes {
        data, err := json.Marshal(cliente)
        if err != nil {
            log.Fatal(err)
        }
        createUpdate(db, "cliente", []byte(strconv.Itoa(cliente.NroCliente)), data)
        resultado, err := readUnique(db, "cliente", []byte(strconv.Itoa(cliente.NroCliente)))
        fmt.Printf("\n%s\n", resultado)
    }
    fmt.Print("\nSE CARGARON DATOS DE CLIENTES\n\n")
    fmt.Printf("\n%s", separador)
    
	//insert a tabla tarjeta
	fmt.Print("\nTABLA TARJETAS\n\n")
	for _, tarjeta := range tarjetas {
        data, err := json.Marshal(tarjeta)
        if err != nil {
            log.Fatal(err)
        }
        createUpdate(db, "tarjeta", []byte(tarjeta.NroTarjeta), data)
        resultado, err := readUnique(db, "tarjeta", []byte(tarjeta.NroTarjeta))
        fmt.Printf("\n%s\n", resultado)
    }
    fmt.Print("\nSE CARGARON DATOS DE TARJETAS\n\n")
    fmt.Printf("\n%s", separador)
    
    //insert a tabla comercio
    fmt.Print("\nTABLA COMERCIO\n\n")
    for _, comercio := range comercios {
        data, err := json.Marshal(comercio)
        if err != nil {
            log.Fatal(err)
        }
        createUpdate(db, "comercio", []byte(strconv.Itoa(comercio.NroComercio)), data)
        resultado, err := readUnique(db, "comercio", []byte(strconv.Itoa(comercio.NroComercio)))
        fmt.Printf("\n%s\n", resultado)
    }
    fmt.Print("\nSE CARGARON DATOS DE COMERCIOS\n\n")
    fmt.Printf("\n%s", separador)
	
	//insert a tabla compra
	fmt.Print("\nTABLA COMPRA\n\n")
	for _, compra := range compras {
        data, err := json.Marshal(compra)
        if err != nil {
            log.Fatal(err)
        }
        createUpdate(db, "compra", []byte(strconv.Itoa(compra.NroOperacion)), data)
        resultado, err := readUnique(db, "compra", []byte(strconv.Itoa(compra.NroOperacion)))
        fmt.Printf("\n%s\n", resultado)
    }
    fmt.Print("\nSE CARGARON DATOS DE COMPRAS\n\n")
    fmt.Printf("\n%s", separador)
}

func createUpdate(db *bolt.DB, bucketName string, key []byte, val []byte) error {
    tx, err := db.Begin(true)
    if err != nil {
        return err
    }
    defer tx.Rollback()
    b, _ := tx.CreateBucketIfNotExists([]byte(bucketName))
    err = b.Put(key, val)
    if err != nil {
        return err
    }
    if err := tx.Commit(); err != nil {
        return err
    }
    return nil
}

func readUnique(db *bolt.DB, bucketName string, key []byte) ([]byte, error) {
    var buf []byte
    err := db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(bucketName))
        buf = b.Get(key)
        return nil
    })
    return buf, err
}
