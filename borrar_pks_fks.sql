--BORRAR PKS
alter table cliente drop constraint clientepk;
alter table tarjeta drop constraint tarjetapk;
alter table comercio drop constraint comerciopk;
alter table compra drop constraint comprapk;
alter table rechazo drop constraint rechazopk;
alter table cierre drop constraint cierrepk;
alter table cabecera drop constraint cabecerapk;
alter table detalle drop constraint detallepk;
alter table alerta drop constraint alertapk;

--BORRAR FKS
alter table tarjeta drop constraint tarjetanroclientefk;
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
