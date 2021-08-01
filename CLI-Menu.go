package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)


type consumo struct {
	nrotarjeta 		string
	codseguridad 	string
	nrocomercio 	int
	monto 			float64
}	


const dbInfo string = "user=postgres host=localhost dbname=tarjetas sslmode=disable"
const separador string = "***********************************************************************************"

var db *sql.DB
var err error
var opcionUsuario int = 0

func main() {
	for opcionUsuario != 9 {
		//Menu programa
		fmt.Print("\nMenu \n")
		fmt.Print("\n------------------------------\n")
		fmt.Print("1 - Crear la base de Tarjetas\n")
		fmt.Print("2 - Crear tablas\n")
		fmt.Print("3 - Crear PKs y FKs\n")
		fmt.Print("4 - Eliminar PKs y FKs\n")
		fmt.Print("5 - Cargar datos a tablas\n")
		fmt.Print("6 - Crear funciones\n")
		fmt.Print("7 - Autorizar compra\n")
		fmt.Print("8 - Generar resumen\n")
		fmt.Print("9 - Salir\n")
		fmt.Print("\n------------------------------\n\n")

		//Menu para seleccion de usuario
		fmt.Print("\nSeleccione una opcion: ")
		fmt.Scanf("%d", &opcionUsuario)

		//Muestra seleccion de usuario
		fmt.Print("\nUsted selecciono: ", opcionUsuario)
		fmt.Print("\n")

		//Seteamos funciones para selecciones del menu
		if opcionUsuario == 1 {
			//Crea DB
			createDataBase()
			db, err = sql.Open("postgres", dbInfo)
			if err != nil {
				log.Fatal(err)
			}
			defer db.Close()
		}

		if opcionUsuario == 2 {
			//Crea Tablas
			createTables()
			if err != nil {
				log.Fatal(err)
			}
			defer db.Close()		
		}

		if opcionUsuario == 3 {
			//Crea PKs
			createPKsFKs()			
		}

		if opcionUsuario == 4 {
			//Elimina PKs
			deletePKsFKs()
		}

		if opcionUsuario == 5 {
			//Poblar Tablas
			insertValues()
		}

		if opcionUsuario == 6 {
			//Crear Funcoines Store Procedures
			init_autorizarCompra()
			init_generarResumen()
			init_alertas()		
		}

		if opcionUsuario == 7 {
			//Autorizar Compras
			autorizar_compra()
			
		}

		if opcionUsuario == 8 {
			//Generar Resumenes
			generar_resumen()
		}

		if opcionUsuario == 9 {
			//Salir
			fmt.Print("\nSalir \n\n")
		}
	}
}

func createDataBase() {
	db, err = sql.Open("postgres", "user=postgres host=localhost dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	//eliminamos db tarjetas si existe
	_, err = db.Exec(`drop database if exists tarjetas`)
	if err != nil {
		log.Fatal(err)
	}
	//creamos db tarjetas
	_, err = db.Exec(`create database tarjetas`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("\nSE CREO BASE CORRECTAMENTE\n\n", separador, "\n\n")

}
func createTables() {
	_, err = db.Exec(`CREATE table cliente (nrocliente int, nombre text, apellido text, domicilio text, telefono char(12));
			CREATE table tarjeta (nrotarjeta char(16), nrocliente int, validadesde char(6), validahasta char(6), codseguridad char(4), limitecompra decimal(8,2), estado char(10));
			CREATE table comercio (nrocomercio int, nombre text, domicilio text, codigopostal char(8), telefono char(12));
			CREATE table compra (nrooperacion serial, nrotarjeta char(16), nrocomercio int, fecha timestamp, monto decimal(7,2), pagado boolean);
			CREATE table rechazo (nrorechazo serial, nrotarjeta char(16), nrocomercio int, fecha timestamp, monto decimal(7,2), motivo text);
			CREATE table cierre (año int, mes int, terminacion int, fechainicio date, fechacierre date, fechavto date);
			CREATE table cabecera (nroresumen serial, nombre text, apellido text, domicilio text, nrotarjeta char(16), desde date, hasta date, vence date, total decimal(8,2));
			CREATE table detalle (nroresumen int, nrolinea int, fecha date, nombrecomercio text, monto decimal (7,2));
			CREATE table alerta (nroalerta serial, nrotarjeta char(16), fecha timestamp, nrorechazo int, codalerta int, descripcion text);
			CREATE table consumo (nrotarjeta char(16), codseguridad char(4), nrocomercio int, monto decimal(7,2));`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("\nSE CREARON LAS TABLAS\n\n",separador,"\n\n")
}

func createPKsFKs() {
	_, err = db.Exec(`alter table cliente add constraint clientepk primary key (nrocliente);
			alter table tarjeta add constraint tarjetapk primary key (nrotarjeta);
			alter table comercio add constraint comerciopk primary key (nrocomercio);
			alter table compra add constraint comprapk primary key (nrooperacion);
			alter table rechazo add constraint rechazopk primary key (nrorechazo);
			alter table cierre add constraint cierrepk primary key (año,mes,terminacion);
			alter table cabecera add constraint cabecerapk primary key (nroresumen);
			alter table detalle add constraint detallepk primary key (nroresumen,nrolinea);
			alter table alerta add constraint alertapk primary key (nroalerta);
			alter table tarjeta add constraint tarjetanroclientefk foreign key (nrocliente) references cliente(nrocliente);				
			alter table compra add constraint compranrotarjetafk foreign key (nrotarjeta) references tarjeta(nrotarjeta);
			alter table compra add constraint compranrocomerciofk foreign key (nrocomercio) references comercio(nrocomercio);
			alter table rechazo add constraint rechazonrocomerciofk foreign key (nrocomercio) references comercio(nrocomercio);
			alter table rechazo add constraint rechazonrotarjetafk foreign key (nrotarjeta) references tarjeta(nrotarjeta);
			alter table cabecera add constraint cabeceranrotarjetafk foreign key (nrotarjeta) references tarjeta(nrotarjeta);
			alter table detalle add constraint detallenroresumenfk foreign key (nroresumen) references cabecera(nroresumen);
			alter table alerta add constraint alertanrotarjetafk foreign key (nrotarjeta) references tarjeta(nrotarjeta);
			alter table alerta add constraint alertanrorechazofk foreign key (nrorechazo) references rechazo(nrorechazo);
			alter table consumo add constraint consumonrotarjetafk foreign key (nrotarjeta) references tarjeta(nrotarjeta);
			alter table consumo add constraint consumonrocomerciofk foreign key (nrocomercio) references comercio(nrocomercio);`)
						
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("\nSE CREARON PKs y FKs \n\n",separador,"\n\n")  
}

func deletePKsFKs() {
	_, err = db.Exec(`alter table tarjeta drop constraint tarjetanroclientefk;
			alter table compra drop constraint compranrotarjetafk;
			alter table compra drop constraint compranrocomerciofk;
			alter table rechazo drop constraint rechazonrocomerciofk;
			alter table rechazo drop constraint rechazonrotarjetafk;
			alter table cabecera drop constraint cabeceranrotarjetafk;
			alter table detalle drop constraint detallenroresumenfk;
			alter table alerta drop constraint alertanrotarjetafk;
			alter table alerta drop constraint alertanrorechazofk;
			alter table consumo drop constraint consumonrotarjetafk; 
			alter table consumo drop constraint consumonrocomerciofk;
			alter table cliente drop constraint clientepk;
			alter table tarjeta drop constraint tarjetapk;
			alter table comercio drop constraint comerciopk;
			alter table compra drop constraint comprapk;
			alter table rechazo drop constraint rechazopk;
			alter table cierre drop constraint cierrepk;
			alter table cabecera drop constraint cabecerapk;
			alter table detalle drop constraint detallepk;
			alter table alerta drop constraint alertapk;`)
			
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("\nSE ELIMINARON PKs y FKs \n\n",separador,"\n\n")

}

func insertValues(){
	_,err = db.Exec (`insert into cliente (nrocliente,nombre,apellido,domicilio,telefono) values ('1','Juan','Riquelme','Bransen 123','541130888931');
			insert into cliente (nrocliente,nombre,apellido,domicilio,telefono) values ('2','Javier','Saviola','Nuñez 234','541130888932');
			insert into cliente (nrocliente,nombre,apellido,domicilio,telefono) values ('3','Juan','Cavaliere','Ricchieri 950','541130888933');
			insert into cliente (nrocliente,nombre,apellido,domicilio,telefono) values ('4','Ignacio','Escapa','Tucuman 800','541130888934');
			insert into cliente (nrocliente,nombre,apellido,domicilio,telefono) values ('5','Nicolas','Correa','Entre Rios 650','541130888935');
			insert into cliente (nrocliente,nombre,apellido,domicilio,telefono) values ('6','Ezequiel','Ruidiaz','Santa fe 650','541130888936');
			insert into cliente (nrocliente,nombre,apellido,domicilio,telefono) values ('7','Dante','De Franco','Paunero 700','541130888937');
			insert into cliente (nrocliente,nombre,apellido,domicilio,telefono) values ('8','Luis','Alcarraz','Ricchieri 838','541130888938');
			insert into cliente (nrocliente,nombre,apellido,domicilio,telefono) values ('9','Diego','Maradona','Moine 1000','541130888939');
			insert into cliente (nrocliente,nombre,apellido,domicilio,telefono) values ('10','Lionel','Messi','Barcelona 222','541130888940');
			insert into cliente (nrocliente,nombre,apellido,domicilio,telefono) values ('11','Carlos','Tevez','Apache 246','541130888941');
			insert into cliente (nrocliente,nombre,apellido,domicilio,telefono) values ('12','Claudio Paul','Caniggia','Italia 90','541130888942');
			insert into cliente (nrocliente,nombre,apellido,domicilio,telefono) values ('13','Gabriel Omar','Batistuta','Francia 98','541130888943');
			insert into cliente (nrocliente,nombre,apellido,domicilio,telefono) values ('14','Gabriela','Sabatini','Estados Unidos 90','541130888944');
			insert into cliente (nrocliente,nombre,apellido,domicilio,telefono) values ('15','Oscar Natalio','Bonavena','Eduardo Madero 470','541130888945');
			insert into cliente (nrocliente,nombre,apellido,domicilio,telefono) values ('16','Martin','Palermo','Corrientes 4193 ','541130888946');
			insert into cliente (nrocliente,nombre,apellido,domicilio,telefono) values ('17','Ikram','Lora', '24 De Septiembre 263','548887273789');
			insert into cliente (nrocliente,nombre,apellido,domicilio,telefono) values ('18','Bruno','Sagarra', 'Galacia 1430','547205527792');
			insert into cliente (nrocliente,nombre,apellido,domicilio,telefono) values ('19','Aroa','Ribes', 'Necochea 62','546382130490');
			insert into cliente (nrocliente,nombre,apellido,domicilio,telefono) values ('20','Braulio','Saldaña', 'Cosme BEccar 274','548652307680');
			insert into comercio (nrocomercio,nombre,domicilio,codigopostal,telefono) values ('1','Supermercado Dia','Sdor Moron 200','1661','46660011');
			insert into comercio (nrocomercio,nombre,domicilio,codigopostal,telefono) values ('2','Supermercado Disco','Sdor Moron 1600','1661','46660022');
			insert into comercio (nrocomercio,nombre,domicilio,codigopostal,telefono) values ('3','Supermercado Coto','Munzon 6200','1662','46660033');
			insert into comercio (nrocomercio,nombre,domicilio,codigopostal,telefono) values ('4','Supermercado Carrefour','Peron 200','1663','46660044');
			insert into comercio (nrocomercio,nombre,domicilio,codigopostal,telefono) values ('5','Supermercado Norte','Pardo 900','1662','46660055');
			insert into comercio (nrocomercio,nombre,domicilio,codigopostal,telefono) values ('6','Carniceria CrisCar','Ricchieri 1100','1661','46660066');
			insert into comercio (nrocomercio,nombre,domicilio,codigopostal,telefono) values ('7','Carniceria Molto','Mitre 1200','1663','46660077');
			insert into comercio (nrocomercio,nombre,domicilio,codigopostal,telefono) values ('8','Carniceria Los Reyes','España 3200','1663','46660088');
			insert into comercio (nrocomercio,nombre,domicilio,codigopostal,telefono) values ('9','Libreria Fox','San Jose 700','1662','46660099');
			insert into comercio (nrocomercio,nombre,domicilio,codigopostal,telefono) values ('10','Libreria La Esquina','Chubut 550','1661','46661100'); 
			insert into comercio (nrocomercio,nombre,domicilio,codigopostal,telefono) values ('11','Verduleria Santa Isabel','Richieri 1054'         ,'1661'   ,'46661314');
			insert into comercio (nrocomercio,nombre,domicilio,codigopostal,telefono) values ('12','MaxiKiosco BV'          ,'Senador Moron 784'     ,'1661'   ,'46661550');
			insert into comercio (nrocomercio,nombre,domicilio,codigopostal,telefono) values ('13','Supermercado EL Oasis'  ,'Alvear Marcelo T De 2323'              ,'92614'  ,'46669852');
			insert into comercio (nrocomercio,nombre,domicilio,codigopostal,telefono) values ('14','Cien Acres'             ,'Libertad 3780'                         ,'67510'  ,'44603929');
			insert into comercio (nrocomercio,nombre,domicilio,codigopostal,telefono) values ('15','La Cesta De Oro'        ,'Mapiu 572'                             ,'72701'  ,'44609381');
			insert into comercio (nrocomercio,nombre,domicilio,codigopostal,telefono) values ('16','El Emporio'             ,'Chacobuco Bat De 1271'                 ,'48823'  ,'46601192');
			insert into comercio (nrocomercio,nombre,domicilio,codigopostal,telefono) values ('17','La Gran Canasta'        ,'Derqui 831'                            ,'02740'  ,'46608281');
			insert into comercio (nrocomercio,nombre,domicilio,codigopostal,telefono) values ('18','Maxi Ahorro'            ,'Gral Paunero 135'                      ,'59013'  ,'46607382');
			insert into comercio (nrocomercio,nombre,domicilio,codigopostal,telefono) values ('19','La Pradera Minimarket'  ,'Maza 3260'                             ,'74120'  ,'44601144');
			insert into comercio (nrocomercio,nombre,domicilio,codigopostal,telefono) values ('20','Tiempo De Ahorro'       ,'Dr.Adolfo Alsina 977'                  ,'38401'  ,'44602858');		
			insert into tarjeta values('4540730040713558', 1, '202001','202401', '356', '15000','vigente');
			insert into tarjeta values('4540730040713578', 1, '201211','201711', '505', '150000','anulada');
			insert into tarjeta values('4540730040713559', 2, '201901','202301', '357', '10000','vigente');
			insert into tarjeta values('4540730040713579', 2, '201612','202012', '506', '160000','anulada');
			insert into tarjeta values('4540730040713560', 3, '201803','202203', '358', '12000','vigente');
			insert into tarjeta values('4540730040713561', 4, '201801','202101', '359', '15000','anulada');
			insert into tarjeta values('4540730040713562', 5, '201802','202102', '360', '20000','anulada');
			insert into tarjeta values('4540730040713563', 6, '201803','202103', '361', '25000','anulada');
			insert into tarjeta values('4540730040713564', 7, '202101','202501', '362', '30000','suspendida');
			insert into tarjeta values('4540730040713565', 8, '202001','202401', '363', '35000','suspendida');
			insert into tarjeta values('4540730040713566', 9, '201712','202112', '364', '40000','suspendida');
			insert into tarjeta values('4540730040713567', 10, '201904','202304', '365', '45000','vigente');
			insert into tarjeta values('4540730040713568', 11, '200707','201107', '366', '50000','anulada');
			insert into tarjeta values('4540730040713569', 12, '202105','202505', '367', '55000','suspendida');
			insert into tarjeta values('4540730040713570', 13, '202104','202504', '368', '60000','suspendida');
			insert into tarjeta values('4540730040713571', 14, '200312','200712', '369', '80000','anulada');
			insert into tarjeta values('4540730040713572', 15, '201906','202406', '370', '90000','vigente');
			insert into tarjeta values('4540730040713573', 16, '201204','201604', '351', '100000','anulada');
			insert into tarjeta values('4540730040713574', 17, '201605','202005', '352', '110000','anulada');
			insert into tarjeta values('4540730040713575', 18, '202105','202505', '353', '115000','vigente');
			insert into tarjeta values('4540730040713576', 19, '202101','202501', '354', '120000','suspendida');
			insert into tarjeta values('4540730040713577', 20, '202012','202512', '355', '125000','vigente');
			insert into cierre values(2021, 1, 0, '2021-01-01', '2021-01-31', '2021-02-11');
			insert into cierre values(2021, 2, 0, '2021-02-01', '2021-02-28', '2021-03-11');
			insert into cierre values(2021, 3, 0, '2021-03-01', '2021-03-31', '2021-04-11');
			insert into cierre values(2021, 4, 0, '2021-04-01', '2021-04-30', '2021-05-11');
			insert into cierre values(2021, 5, 0, '2021-05-01', '2021-05-31', '2021-06-11');
			insert into cierre values(2021, 6, 0, '2021-06-01', '2021-06-30', '2021-07-11');
			insert into cierre values(2021, 7, 0, '2021-07-01', '2021-07-31', '2021-08-11');
			insert into cierre values(2021, 8, 0, '2021-08-01', '2021-08-31', '2021-09-11');
			insert into cierre values(2021, 9, 0, '2021-09-01', '2021-09-30', '2021-10-11');
			insert into cierre values(2021, 10, 0, '2021-10-01', '2021-10-31', '2021-11-11');
			insert into cierre values(2021, 11, 0, '2021-11-01', '2021-11-30', '2021-12-11');
			insert into cierre values(2021, 12, 0, '2021-12-01', '2021-12-31', '2022-01-11');
			insert into cierre values(2021, 1, 1, '2021-01-02', '2021-02-01', '2021-02-12');
			insert into cierre values(2021, 2, 1, '2021-02-02', '2021-03-01', '2021-03-12');
			insert into cierre values(2021, 3, 1, '2021-03-02', '2021-04-01', '2021-04-12');
			insert into cierre values(2021, 4, 1, '2021-04-02', '2021-05-01', '2021-05-12');
			insert into cierre values(2021, 5, 1, '2021-05-02', '2021-06-01', '2021-06-12');
			insert into cierre values(2021, 6, 1, '2021-06-02', '2021-07-01', '2021-07-12');
			insert into cierre values(2021, 7, 1, '2021-07-02', '2021-08-01', '2021-08-12');
			insert into cierre values(2021, 8, 1, '2021-08-02', '2021-09-01', '2021-09-12');
			insert into cierre values(2021, 9, 1, '2021-09-02', '2021-10-01', '2021-10-12');
			insert into cierre values(2021, 10, 1, '2021-10-02', '2021-11-01', '2021-11-12');
			insert into cierre values(2021, 11, 1, '2021-11-02', '2021-12-01', '2021-12-12');
			insert into cierre values(2021, 12, 1, '2021-12-02', '2022-01-01', '2022-01-12');
			insert into cierre values(2021, 1, 2, '2021-01-03', '2021-02-02', '2021-02-13');
			insert into cierre values(2021, 2, 2, '2021-02-03', '2021-03-02', '2021-03-13');
			insert into cierre values(2021, 3, 2, '2021-03-03', '2021-04-02', '2021-04-13');
			insert into cierre values(2021, 4, 2, '2021-04-03', '2021-05-02', '2021-05-13');
			insert into cierre values(2021, 5, 2, '2021-05-03', '2021-06-02', '2021-06-13');
			insert into cierre values(2021, 6, 2, '2021-06-03', '2021-07-02', '2021-07-13');
			insert into cierre values(2021, 7, 2, '2021-07-03', '2021-08-02', '2021-08-13');
			insert into cierre values(2021, 8, 2, '2021-08-03', '2021-09-02', '2021-09-13');
			insert into cierre values(2021, 9, 2, '2021-09-03', '2021-10-02', '2021-10-13');
			insert into cierre values(2021, 10, 2, '2021-10-03', '2021-11-02', '2021-11-13');
			insert into cierre values(2021, 11, 2, '2021-11-03', '2021-12-02', '2021-12-13');
			insert into cierre values(2021, 12, 2, '2021-12-03', '2022-01-02', '2022-01-13');
			insert into cierre values(2021, 1, 3, '2021-01-04', '2021-02-03', '2021-02-14');
			insert into cierre values(2021, 2, 3, '2021-02-04', '2021-03-03', '2021-03-14');
			insert into cierre values(2021, 3, 3, '2021-03-04', '2021-04-03', '2021-04-14');
			insert into cierre values(2021, 4, 3, '2021-04-04', '2021-05-03', '2021-05-14');
			insert into cierre values(2021, 5, 3, '2021-05-04', '2021-06-03', '2021-06-14');
			insert into cierre values(2021, 6, 3, '2021-06-04', '2021-07-03', '2021-07-14');
			insert into cierre values(2021, 7, 3, '2021-07-04', '2021-08-03', '2021-08-14');
			insert into cierre values(2021, 8, 3, '2021-08-04', '2021-09-03', '2021-09-14');
			insert into cierre values(2021, 9, 3, '2021-09-04', '2021-10-03', '2021-10-14');
			insert into cierre values(2021, 10, 3, '2021-10-04', '2021-11-03', '2021-11-14');
			insert into cierre values(2021, 11, 3, '2021-11-04', '2021-12-03', '2021-12-14');
			insert into cierre values(2021, 12, 3, '2021-12-04', '2022-01-03', '2022-01-14');
			insert into cierre values(2021, 1, 4, '2021-01-05', '2021-02-04', '2021-02-15');
			insert into cierre values(2021, 2, 4, '2021-02-05', '2021-03-04', '2021-03-15');
			insert into cierre values(2021, 3, 4, '2021-03-05', '2021-04-04', '2021-04-15');
			insert into cierre values(2021, 4, 4, '2021-04-05', '2021-05-04', '2021-05-15');
			insert into cierre values(2021, 5, 4, '2021-05-05', '2021-06-04', '2021-06-15');
			insert into cierre values(2021, 6, 4, '2021-06-05', '2021-07-04', '2021-07-15');
			insert into cierre values(2021, 7, 4, '2021-07-05', '2021-08-04', '2021-08-15');
			insert into cierre values(2021, 8, 4, '2021-08-05', '2021-09-04', '2021-09-15');
			insert into cierre values(2021, 9, 4, '2021-09-05', '2021-10-04', '2021-10-15');
			insert into cierre values(2021, 10, 4, '2021-10-05', '2021-11-04', '2021-11-15');
			insert into cierre values(2021, 11, 4, '2021-11-05', '2021-12-04', '2021-12-15');
			insert into cierre values(2021, 12, 4, '2021-12-05', '2022-01-04', '2022-01-15');
			insert into cierre values(2021, 1, 5, '2021-01-06', '2021-02-05', '2021-02-16');
			insert into cierre values(2021, 2, 5, '2021-02-06', '2021-03-05', '2021-03-16');
			insert into cierre values(2021, 3, 5, '2021-03-06', '2021-04-05', '2021-04-16');
			insert into cierre values(2021, 4, 5, '2021-04-06', '2021-05-05', '2021-05-16');
			insert into cierre values(2021, 5, 5, '2021-05-06', '2021-06-05', '2021-06-16');
			insert into cierre values(2021, 6, 5, '2021-06-06', '2021-07-05', '2021-07-16');
			insert into cierre values(2021, 7, 5, '2021-07-06', '2021-08-05', '2021-08-16');
			insert into cierre values(2021, 8, 5, '2021-08-06', '2021-09-05', '2021-09-16');
			insert into cierre values(2021, 9, 5, '2021-09-06', '2021-10-05', '2021-10-16');
			insert into cierre values(2021, 10, 5, '2021-10-06', '2021-11-05', '2021-11-16');
			insert into cierre values(2021, 11, 5, '2021-11-06', '2021-12-05', '2021-12-16');
			insert into cierre values(2021, 12, 5, '2021-12-06', '2022-01-05', '2022-01-16');
			insert into cierre values(2021, 1, 6, '2021-01-07', '2021-02-06', '2021-02-17');
			insert into cierre values(2021, 2, 6, '2021-02-07', '2021-03-06', '2021-03-17');
			insert into cierre values(2021, 3, 6, '2021-03-07', '2021-04-06', '2021-04-17');
			insert into cierre values(2021, 4, 6, '2021-04-07', '2021-05-06', '2021-05-17');
			insert into cierre values(2021, 5, 6, '2021-05-07', '2021-06-06', '2021-06-17');
			insert into cierre values(2021, 6, 6, '2021-06-07', '2021-07-06', '2021-07-17');
			insert into cierre values(2021, 7, 6, '2021-07-07', '2021-08-06', '2021-08-17');
			insert into cierre values(2021, 8, 6, '2021-08-07', '2021-09-06', '2021-09-17');
			insert into cierre values(2021, 9, 6, '2021-09-07', '2021-10-06', '2021-10-17');
			insert into cierre values(2021, 10, 6, '2021-10-07', '2021-11-06', '2021-11-17');
			insert into cierre values(2021, 11, 6, '2021-11-07', '2021-12-06', '2021-12-17');
			insert into cierre values(2021, 12, 6, '2021-12-07', '2022-01-06', '2022-01-17');
			insert into cierre values(2021, 1, 7, '2021-01-08', '2021-02-07', '2021-02-18');
			insert into cierre values(2021, 2, 7, '2021-02-08', '2021-03-07', '2021-03-18');
			insert into cierre values(2021, 3, 7, '2021-03-08', '2021-04-07', '2021-04-18');
			insert into cierre values(2021, 4, 7, '2021-04-08', '2021-05-07', '2021-05-18');
			insert into cierre values(2021, 5, 7, '2021-05-08', '2021-06-07', '2021-06-18');
			insert into cierre values(2021, 6, 7, '2021-06-08', '2021-07-07', '2021-07-18');
			insert into cierre values(2021, 7, 7, '2021-07-08', '2021-08-07', '2021-08-18');
			insert into cierre values(2021, 8, 7, '2021-08-08', '2021-09-07', '2021-09-18');
			insert into cierre values(2021, 9, 7, '2021-09-08', '2021-10-07', '2021-10-18');
			insert into cierre values(2021, 10, 7, '2021-10-08', '2021-11-07', '2021-11-18');
			insert into cierre values(2021, 11, 7, '2021-11-08', '2021-12-07', '2021-12-18');
			insert into cierre values(2021, 12, 7, '2021-12-08', '2022-01-07', '2022-01-18');
			insert into cierre values(2021, 1, 8, '2021-01-09', '2021-02-08', '2021-02-19');
			insert into cierre values(2021, 2, 8, '2021-02-09', '2021-03-08', '2021-03-19');
			insert into cierre values(2021, 3, 8, '2021-03-09', '2021-04-08', '2021-04-19');
			insert into cierre values(2021, 4, 8, '2021-04-09', '2021-05-08', '2021-05-19');
			insert into cierre values(2021, 5, 8, '2021-05-09', '2021-06-08', '2021-06-19');
			insert into cierre values(2021, 6, 8, '2021-06-09', '2021-07-08', '2021-07-19');
			insert into cierre values(2021, 7, 8, '2021-07-09', '2021-08-08', '2021-08-19');
			insert into cierre values(2021, 8, 8, '2021-08-09', '2021-09-08', '2021-09-19');
			insert into cierre values(2021, 9, 8, '2021-09-09', '2021-10-08', '2021-10-19');
			insert into cierre values(2021, 10, 8, '2021-10-09', '2021-11-08', '2021-11-19');
			insert into cierre values(2021, 11, 8, '2021-11-09', '2021-12-08', '2021-12-19');
			insert into cierre values(2021, 12, 8, '2021-12-09', '2022-01-08', '2022-01-19');
			insert into cierre values(2021, 1, 9, '2021-01-10', '2021-02-09', '2021-02-20');
			insert into cierre values(2021, 2, 9, '2021-02-10', '2021-03-09', '2021-03-20');
			insert into cierre values(2021, 3, 9, '2021-03-10', '2021-04-09', '2021-04-20');
			insert into cierre values(2021, 4, 9, '2021-04-10', '2021-05-09', '2021-05-20');
			insert into cierre values(2021, 5, 9, '2021-05-10', '2021-06-09', '2021-06-20');
			insert into cierre values(2021, 6, 9, '2021-06-10', '2021-07-09', '2021-07-20');
			insert into cierre values(2021, 7, 9, '2021-07-10', '2021-08-09', '2021-08-20');
			insert into cierre values(2021, 8, 9, '2021-08-10', '2021-09-09', '2021-09-20');
			insert into cierre values(2021, 9, 9, '2021-09-10', '2021-10-09', '2021-10-20');
			insert into cierre values(2021, 10, 9, '2021-10-10', '2021-11-09', '2021-11-20');
			insert into cierre values(2021, 11, 9, '2021-11-10', '2021-12-09', '2021-12-20');
			insert into cierre values(2021, 12, 9, '2021-12-10', '2022-01-09', '2022-01-20');			
			insert into consumo (nrotarjeta,codseguridad,nrocomercio,monto) values ('4540730040713559','357',1,'5000');
			insert into consumo (nrotarjeta,codseguridad,nrocomercio,monto) values ('4540730040713560','358',8,'6000');
			insert into consumo (nrotarjeta,codseguridad,nrocomercio,monto) values ('4540730040713567','365',2,'4500');
			insert into consumo (nrotarjeta,codseguridad,nrocomercio,monto) values ('4540730040713572','370',6,'3200');
			insert into consumo (nrotarjeta,codseguridad,nrocomercio,monto) values ('4540730040713577','355',7,'5500');
			insert into consumo (nrotarjeta,codseguridad,nrocomercio,monto) values ('4540730040713558','356',3,'1000');
			insert into consumo (nrotarjeta,codseguridad,nrocomercio,monto) values ('4540730040713560','358',11,'4000');
			insert into consumo (nrotarjeta,codseguridad,nrocomercio,monto) values ('4540730040713579','506',1,'10000');
			insert into consumo (nrotarjeta,codseguridad,nrocomercio,monto) values ('4540730040713575','353',9,'10000');
			insert into consumo (nrotarjeta,codseguridad,nrocomercio,monto) values ('4540730040713578','505',3,'7500');
			insert into consumo (nrotarjeta,codseguridad,nrocomercio,monto) values ('4540730040713578','505',4,'2500');
			insert into consumo (nrotarjeta,codseguridad,nrocomercio,monto) values ('4540730040713562','360',10,'10000');
			insert into consumo (nrotarjeta,codseguridad,nrocomercio,monto) values ('4540730040713558','356',6,'16000');
			insert into consumo (nrotarjeta,codseguridad,nrocomercio,monto) values ('4540730040713559','357',2,'3000');
			insert into consumo (nrotarjeta,codseguridad,nrocomercio,monto) values ('4540730040713565','363',19,'10000');
			insert into consumo (nrotarjeta,codseguridad,nrocomercio,monto) values ('4540730040713577','354',14,'1500');
			insert into consumo (nrotarjeta,codseguridad,nrocomercio,monto) values ('4540730040713558','156',4,'2000');
			insert into consumo (nrotarjeta,codseguridad,nrocomercio,monto) values ('4540730040713572','370',1,'91000');
			insert into consumo (nrotarjeta,codseguridad,nrocomercio,monto) values ('4540730040713572','370',2,'91000');`)

	if err != nil {
		log.Fatal(err)
}

	fmt.Print("\nSE CARGARON DATOS CORRECTAMENTE\n\n")
}

func init_autorizarCompra() {
	_, err = db.Exec(`create or replace function autorizar_compra(numtarjeta tarjeta.nrotarjeta%type,
																  codigo tarjeta.codseguridad%type,
																  numcomercio comercio.nrocomercio%type,
															      monto_abonado compra.monto%type) returns boolean as $$

			declare
					tarjeta_tmp record;
					compras_pend_pago compra.monto%type;
					

			begin
					select * into tarjeta_tmp from tarjeta t where numtarjeta = t.nrotarjeta;
					compras_pend_pago := (select sum(monto) from compra c where c.nrotarjeta = numtarjeta and c.pagado = false);
					if not found then			
						insert into rechazo(nrotarjeta, nrocomercio, fecha, monto, motivo) values (numtarjeta, numcomercio, current_timestamp(0), monto_abonado, 'tarjeta no vàlida o no vigente');
						return false;
					
					elseif tarjeta_tmp.codseguridad != codigo then		
						insert into rechazo(nrotarjeta, nrocomercio, fecha, monto, motivo) values (numtarjeta, numcomercio, current_timestamp(0), monto_abonado, 'còdigo de seguridad invàlido');
						return false;
							
					elseif tarjeta_tmp.limitecompra < (compras_pend_pago + monto_abonado) then
						insert into rechazo(nrotarjeta, nrocomercio, fecha, monto, motivo) values (numtarjeta, numcomercio, current_timestamp(0), monto_abonado, 'supera lìmite de tarjeta');
						return false;
					
					elseif tarjeta_tmp.estado = 'anulada' then		
						insert into rechazo(nrotarjeta, nrocomercio, fecha, monto, motivo) values (numtarjeta, numcomercio, current_timestamp(0), monto_abonado, 'plazo de vigencia expirado');
						return false;
						
					elseif tarjeta_tmp.estado = 'suspendida' then 		
						insert into rechazo(nrotarjeta, nrocomercio, fecha, monto, motivo) values (numtarjeta, numcomercio, current_timestamp(0), monto_abonado, 'la tarjeta se encuentra suspendida');
						return false;
					
					else		
						insert into compra(nrotarjeta, nrocomercio, fecha, monto, pagado) values (numtarjeta, numcomercio, current_timestamp(0), monto_abonado, false);
						return true;	
					end if;				
			end; 
			$$ language plpgsql;`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("\nFUNCION autorizar_compra CREADA CORRECTAMENTE\n\n")			
}

func init_generarResumen() {
	_, err = db.Exec(`create or replace function generar_resumen(ncliente cliente.nrocliente%type,año_tmp int,mes_tmp int) returns void as $$
			declare 
				clienteAux record;
				compraAux record;
				tarjetaAux record;
				cierreAux record;
				totalAux cabecera.total%type;
				nroResumenAux cabecera.nroresumen%type;
				nombreComercio comercio.nombre%type;
				cont int :=1;
			
			begin
				select * into clienteAux from cliente where nrocliente= ncliente;
					if not found then
						raise 'Cliente % no existe', ncliente;
					end if; 
					
				for tarjetaAux in select * from tarjeta where nrocliente = ncliente loop
					totalAux :=0;
					select * into cierreAux from cierre where cierre.año = año_tmp and cierre.mes = mes_tmp and cierre.terminacion = substring (tarjetaAux.nrotarjeta,16,1) ::int;
					
					insert into cabecera(nombre,apellido,domicilio,nrotarjeta,desde,hasta,vence) values (clienteAux.nombre,clienteAux.apellido,clienteAux.domicilio,tarjetaAux.nrotarjeta,cierreAux.fechaInicio,cierreAux.fechaCierre,cierreAux.fechaVto);
					
					select into nroresumenAux nroresumen from cabecera where nrotarjeta = tarjetaAux.nrotarjeta and desde = cierreAux.fechaInicio and hasta = cierreAux.fechaCierre;
					
					for compraAux in select * from compra where nrotarjeta = tarjetaAux.nrotarjeta and fecha::date>=(cierreAux.fechaInicio)::date and fecha::date<=(cierreAux.fechaCierre)::date and pagado = false loop
					
						nombreComercio := (select nombre from comercio where nrocomercio = compraAux.nrocomercio);
						insert into detalle values (nroResumenAux,cont,compraAux.fecha,nombreComercio,compraAux.monto);
						totalAux := totalAux + compraAux.monto;
						cont := cont + 1;
						update compra set pagado = true where nrooperacion = compraAux.nrooperacion;
					end loop; 
					
					update cabecera set total = totalAux where nrotarjeta = tarjetaAux.nrotarjeta and desde = cierreAux.fechaInicio and hasta = cierreAux.fechaCierre;
				end loop;
			end;
			$$ language plpgsql;`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("\nFUNCION generar_resumen CREADA CORRECTAMENTE\n\n")			
}

func init_alertas() {
	_, err = db.Exec(`create or replace function alerta_compras() returns trigger as $$
			begin
				
				perform * from compra c where c.nrotarjeta = new.nrotarjeta and 
				c.nrooperacion != new.nrooperacion and 
				c.fecha >= new.fecha - (1 * interval '1 minute') and 
				c.nrocomercio != new.nrocomercio and 
				(select codigopostal from comercio where nrocomercio = c.nrocomercio) = (select codigopostal from comercio where nrocomercio = new.nrocomercio);
				
				if found then 
					insert into alerta (nrotarjeta,fecha,codalerta,descripcion) 
					values (new.nrotarjeta,current_timestamp(0),1,'Se han registrado dos compras en menos de un minuto en comercios distintos');
				end if;
				
				perform * from compra c where c.nrotarjeta = new.nrotarjeta and 
				c.nrooperacion != new.nrooperacion and 
				c.fecha >= new.fecha - (5 * interval '1 minute') and 
				c.nrocomercio != new.nrocomercio and 
				(select codigopostal from comercio where nrocomercio = c.nrocomercio) != (select codigopostal from comercio where nrocomercio = new.nrocomercio);
				
				if found then 
					insert into alerta (nrotarjeta,fecha,codalerta,descripcion) 
					values (new.nrotarjeta,current_timestamp(0),5,'Se han registrado dos compras en menos de 5 minutos en comercios muy alejados');
				end if;
				return new;
			end;

		$$ language plpgsql;

		create or replace function alerta_rechazos() returns trigger as $$
			begin 
				insert into alerta (nrotarjeta,fecha,nrorechazo,codalerta,descripcion) 
				values (new.nrotarjeta,current_timestamp(0),new.nrorechazo,0,concat('Se han rechazado su compra porque ', new.motivo));
				
				perform * from rechazo r where r.nrotarjeta = new.nrotarjeta and 
				r.nrorechazo != new.nrorechazo and 
				extract(day from r.fecha) = extract(day from new.fecha) and 
				extract(month from r.fecha) = extract(month from new.fecha) and 
				extract(year from r.fecha) = extract(year from new.fecha) and 
				r.motivo = 'supera lìmite de tarjeta' and 
				new.motivo = 'supera lìmite de tarjeta';
				
				if found then
					update tarjeta set estado = 'suspendida' where nrotarjeta = new.nrotarjeta;
					
					insert into alerta (nrotarjeta,fecha,nrorechazo,codalerta,descripcion) 
					values (new.nrotarjeta,current_timestamp(0),new.nrorechazo,32,'Se ha suspendido la tarjeta preventivamente');
				end if;
				return new;
			end;

			$$ language plpgsql;

			create trigger compras_tgr after insert on compra for each row execute procedure alerta_compras();

			create trigger rechazos_tgr after insert on rechazo for each row execute procedure alerta_rechazos();`)
			
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("\nFUNCIONES alerta_compras() Y alerta_rechazos() CREADAS CORRECTAMENTE\n\n")	
}		

func autorizar_compra() {
	
	row, err := db.Query(`select * from consumo`)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var c consumo
	//cicla en registros de consumo 
	for row.Next(){
		
		if err := row.Scan(&c.nrotarjeta, &c.codseguridad, &c.nrocomercio, &c.monto); err != nil {
			log.Fatal(err)	
		}
		var autorizaCompra bool
		//query para llamar a funcion con datos de registro de consumo
		sqlBusqueda := `select autorizar_compra($1, $2, $3, $4)`
		err := db.QueryRow(sqlBusqueda, c.nrotarjeta, c.codseguridad, c.nrocomercio, c.monto).Scan(&autorizaCompra)
		if err != nil {
			log.Fatal(err)
		}
		if autorizaCompra {
			fmt.Print ("\nLA COMRA CON LA TARJETA NUMERO: ", c.nrotarjeta, " FUE AUTORIZADA\n\n", separador, "\n")	
		} else {
			fmt.Print ("\nLA COMRA CON LA TARJETA NUMERO: ", c.nrotarjeta, " FUE RECHAZADA\n\n", separador, "\n")	
		}		
	}
	fmt.Print("\nNO QUEDAN CONSUMOS POR AUTORIZAR\n\n", separador, "\n")					
}	
	
func generar_resumen() {
	//Generamos resumenes para compras de usuarios con tarjeta vigente
	_, err = db.Query(`select generar_resumen(2, 2021, 6);
			select generar_resumen(3, 2021, 6);
			select generar_resumen(1, 2021, 6);`)
			//select generar_resumen(10, 2021, 6);
			//select generar_resumen(15, 2021, 6);
			//select generar_resumen(20, 2021, 6);
			//select generar_resumen(1, 2021, 6);
			//select generar_resumen(18, 2021, 6);`)
	if err != nil {
			log.Fatal(err)
	}	
	fmt.Print("\nRESUMENES GENERADOS CORRECTAMENTE\n\n", separador, "\n")	
}	

