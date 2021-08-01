create or replace function alerta_compras() returns trigger as $$
	begin
		--verifica si se registran dos compras con la misma tarjeta en comercios con mismo codigo postal en lapso menor a un minuto
		perform * from compra c where c.nrotarjeta = new.nrotarjeta and 
		c.nrooperacion != new.nrooperacion and 
		c.fecha >= new.fecha - (1 * interval '1 minute') and 
		c.nrocomercio != new.nrocomercio and 
		(select codigopostal from comercio where nrocomercio = c.nrocomercio) = (select codigopostal from comercio where nrocomercio = new.nrocomercio);
		
		if found then 
			insert into alerta (nrotarjeta,fecha,codalerta,descripcion) 
			values (new.nrotarjeta,current_timestamp(0),1,'Se han registrado dos compras en menos de un minuto en comercios distintos');
		end if;
		
		--verifica si se registran dos compras con la misma tarjeta en comercios con distintos codigo postal en lapso menor a 5 minutos
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
		--genera nueva alerta ante un nuevo rechazo 
		insert into alerta (nrotarjeta,fecha,nrorechazo,codalerta,descripcion) 
		values (new.nrotarjeta,current_timestamp(0),new.nrorechazo,0,concat('Se han rechazado su compra porque ', new.motivo));
		
		--verifica si se registran dos rechazos con la misma tarjeta el mismo dia y ambos por exceso en el limite de compra
		perform * from rechazo r where r.nrotarjeta = new.nrotarjeta and 
		r.nrorechazo != new.nrorechazo and 
		extract(day from r.fecha) = extract(day from new.fecha) and 
		extract(month from r.fecha) = extract(month from new.fecha) and 
		extract(year from r.fecha) = extract(year from new.fecha) and 
		r.motivo = 'supera lìmite de tarjeta' and 
		new.motivo = 'supera lìmite de tarjeta';
		
		if found then
			--actualiza el estado de la tarjeta a suspendida
			update tarjeta set estado = 'suspendida' where nrotarjeta = new.nrotarjeta;
			--genera alerta por suspension de tarjeta
			insert into alerta (nrotarjeta,fecha,nrorechazo,codalerta,descripcion) 
			values (new.nrotarjeta,current_timestamp(0),new.nrorechazo,32,'Se ha suspendido la tarjeta preventivamente');
		end if;
		return new;
	end;

$$ language plpgsql;

create trigger compras_tgr after insert on compra for each row execute procedure alerta_compras();
create trigger rechazos_tgr after insert on rechazo for each row execute procedure alerta_rechazos();
