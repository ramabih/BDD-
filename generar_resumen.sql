create or replace function generar_resumen(ncliente cliente.nrocliente%type, año_tmp int, mes_tmp int) returns void as $$
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
		--Verifica que exista el ncliente en la tabla cliente
		select * into clienteAux from cliente where nrocliente= ncliente;
			if not found then
				raise 'Cliente % no existe ', ncliente;
			end if; 
		--Busca en la tabla tarjeta la tarjeta con nrocliente que recibe por paràmetro	
		for tarjetaAux in select * from tarjeta where nrocliente = ncliente loop
			--acumula monto de compras
			totalAux :=0;
			--select en la tabla cierre para la tarjeta(ultimo nro de tarjeta) del nrocliente pasado por parametro en el periodo pasado por parametro
			select * into cierreAux from cierre where cierre.año = año_tmp and cierre.mes = mes_tmp and cierre.terminacion = substring (tarjetaAux.nrotarjeta,16,1) ::int;
			--inserta en tabla cabecera los datos del cliente 
			insert into cabecera(nombre,apellido,domicilio,nrotarjeta,desde,hasta,vence) values (clienteAux.nombre,clienteAux.apellido,clienteAux.domicilio,tarjetaAux.nrotarjeta,cierreAux.fechaInicio,cierreAux.fechaCierre,cierreAux.fechaVto);
			
			select into nroresumenAux nroresumen from cabecera where nrotarjeta = tarjetaAux.nrotarjeta and desde = cierreAux.fechaInicio and hasta = cierreAux.fechaCierre;
			--cicla en tabla compra para insertar las compras del cliente en tabla detalle
			for compraAux in select * from compra where nrotarjeta = tarjetaAux.nrotarjeta and fecha::date>=(cierreAux.fechaInicio)::date and fecha::date<=(cierreAux.fechaCierre)::date and pagado = false loop
			
				nombreComercio := (select nombre from comercio where nrocomercio = compraAux.nrocomercio);
				insert into detalle values (nroResumenAux,cont,compraAux.fecha,nombreComercio,compraAux.monto);
				totalAux := totalAux + compraAux.monto;
				cont := cont + 1;
				--actualiza estado de pago en la compra
				update compra set pagado = true where nrooperacion = compraAux.nrooperacion;
			end loop; 
			--actualiza monto total en tabla cabecera
			update cabecera set total = totalAux where nrotarjeta = tarjetaAux.nrotarjeta and desde = cierreAux.fechaInicio and hasta = cierreAux.fechaCierre;
		end loop;
	end;
	
$$ language plpgsql;
			
			
