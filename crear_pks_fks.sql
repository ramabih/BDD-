--PKs
alter table cliente add constraint clientepk primary key (nrocliente);
alter table tarjeta add constraint tarjetapk primary key (nrotarjeta);
alter table comercio add constraint comerciopk primary key (nrocomercio);
alter table compra add constraint comprapk primary key (nrooperacion);
alter table rechazo add constraint rechazopk primary key (nrorechazo);
alter table cierre add constraint cierrepk primary key (año,mes,terminacion);
alter table cabecera add constraint cabecerapk primary key (nroresumen);
alter table detalle add constraint detallepk primary key (nroresumen,nrolinea);
alter table alerta add constraint alertapk primary key (nroalerta);

--FKs
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
alter table consumo add constraint consumonrocomerciofk foreign key (nrocomercio) references comercio(nrocomercio); 
