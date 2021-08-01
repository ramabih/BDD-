create or replace function autorizar_compra(numtarjeta tarjeta.nrotarjeta%type,
											codigo tarjeta.codseguridad%type,
											numcomercio comercio.nrocomercio%type,
										    monto_abonado compra.monto%type) returns boolean as $$

declare
		
		tarjeta_tmp record;
		compras_pend_pago compra.monto%type;

begin
		select * into tarjeta_tmp from tarjeta t where numtarjeta = t.nrotarjeta;
		compras_pend_pago := (select sum(monto) from compra c where c.nrotarjeta = numtarjeta and c.pagado = false);
		
		--El nrotarjeta no correponde a tarjeta vigente o valida, genera rechazo de la compra
		if not found then			
			insert into rechazo (nrotarjeta, nrocomercio, fecha, monto, motivo) values (numtarjeta, numcomercio, current_timestamp, monto_abonado, 'tarjeta no vàlida o no vigente');
		    return false;
		
		--El codseguridad de la tarjeta no es valido, genera rechazo de la compra
		elseif tarjeta_tmp.codseguridad != codigo then		
			insert into rechazo (nrotarjeta, nrocomercio, fecha, monto, motivo) values (numtarjeta, numcomercio, current_timestamp, monto_abonado, 'còdigo de seguridad invàlido');
			return false;
			
		--Las compras pendientes+monto abonado supera el limite de la tarjeta, genera rechazo de la compra	
		elseif tarjeta_tmp.limitecompra <= (compras_pend_pago + monto_abonado) then
			insert into rechazo (nrotarjeta, nrocomercio, fecha, monto, motivo) values (numtarjeta, numcomercio, current_timestamp, monto_abonado, 'supera lìmite de tarjeta');
			return false;
		
		--La tarjeta esta vencida, genera rechazo de la compra	
		elseif tarjeta_tmp.estado = 'anulada' then		
			insert into rechazo (nrotarjeta, nrocomercio, fecha, monto, motivo) values (numtarjeta, numcomercio, current_timestamp, monto_abonado, 'plazo de vigencia expirado');
			return false;
		
		--La tareta se encuentra suspendida, genera rechazo de la compra	
		elseif tarjeta_tmp.estado = 'suspendida' then 		
			insert into rechazo (nrotarjeta, nrocomercio, fecha, monto, motivo) values (numtarjeta, numcomercio, current_timestamp, monto_abonado, 'la tarjeta se encuentra suspendida');
			return false;
		
		--El consumo es valido, genera nueva compra	
		else		
			insert into compra (nrotarjeta, nrocomercio, fecha, monto, pagado) values (numtarjeta, numcomercio, current_timestamp, monto_abonado, false);
			return true;
			
		end if;	
			
end; $$ language plpgsql;


			
			
			
			
			
			
			
			
			
			
			
			
