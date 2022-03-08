<?php include 'DGME/estatico/encabezado.php';
include 'ControlMigratorio/comun/TodoControlMigratorio.php';



if ($_SESSION['idUnidad'] == 29) {
    if (isset($_REQUEST['opcion'])) {
        $_SESSION['opcion'] = $_REQUEST['opcion'];
    } else {
        $_SESSION['opcion'];
    }

    if (isset($_REQUEST['tipo'])) {
        $_SESSION['tipo'] = $_REQUEST['tipo'];
    } else {
        $_SESSION['tipo'];
    }
}

if ($_SESSION['idUnidad'] == 28) {
    if (isset($_REQUEST['opcion'])) {
        $_SESSION['opcion'] = $_REQUEST['opcion'];
    } else {
        $_SESSION['opcion'];
    }
}

$tablaAtributos = "width=80% style='font-family:Arial; font-size:12px' align=center";
TituloAeropuerto($_SESSION['idUnidad'], $_SESSION['opcion'], $_SESSION['tipo']);

c::inicioCentrado();
?>
    <!--<iframe src="/ControlMigratorio/Aeropuertos/teclado_lector.php" height="100" width="600" frameborder="0" ></iframe>-->
<script type="text/javascript">
	var ws;
    
    function conectar(){
    	document.getElementById('estado').value = 'Conectando...';
		ws = new WebSocket('ws://localhost:8585/LectorOCR');
		ws.onopen = function() { document.getElementById('estado').value = 'Conexion Establecida.'; };
		ws.onmessage = function (evento) { 
			var mensaje = evento.data;
			if (mensaje.includes('LAMINAJPG64') ){
				var arrayDeCadenas = mensaje.split('(|)');
				document.getElementById('lamina').src = arrayDeCadenas[1];
			}
			else if (mensaje.includes('FACEJPG64') ){
				var arrayDeCadenas = mensaje.split('(|)');
				document.getElementById('face').src = arrayDeCadenas[1];
			}
			else if (mensaje.includes('OCR') ){
				var arrayDeCadenas = mensaje.split('(|)');
				document.getElementById('txtTerminal').value = arrayDeCadenas[1];
			}
			else{
				document.getElementById('estado').value = mensaje;
			}
		};
		ws.onclose = function() { 
			document.getElementById('estado').value = 'Desconectado.';
			window.location.assign('appurl://argumentos');
			self.focus();
			setTimeout(function() { conectar(); }, 1000);
		};
	}

	function sendws(){
		$("#txtTerminal").focus();
		var e = jQuery.Event("keydown");
		e.which = 13; // # Some key code value
		e.keyCode = 13;
		$("#txtTerminal").trigger(e);
		$("#txtTerminal").focus();
		alert("Enviado");
	}

    function camara(){
        //abrirFrameGenerico("camara.php",null,395,620);
        //window.open("camara.php" , "ventana1" , "width=620,height=395,scrollbars=NO")
        window.open ("http://localhost:8080/camara/","mywindow", "width=420,height=630");
    }

    // Create IE + others compatible event handler
          var eventMethod = window.addEventListener ? "addEventListener" : "attachEvent";
          var eventer = window[eventMethod];
          var messageEvent = eventMethod == "attachEvent" ? "onmessage" : "message";

          // Listen to message from child window
          eventer(messageEvent,function(e) {
            console.log('origin: ', e.origin)
            
            // Check if origin is proper
            //if( e.origin != 'http://localhost:8080' ){ return }
            //if( e.origin != '*' ){ return }

            //console.log('parent received message!: ', e.data);

            //alert(e.data);
            document.getElementById('face').src = e.data;

          }, false);



    /*var eventMethod = window.addEventListener ? "addEventListener" : "attachEvent";
    var eventer = window[eventMethod];
    var messageEvent = eventMethod === "attachEvent" ? "onmessage" : "message";
    eventer(messageEvent, function (e) {
        document.getElementById('face').src = e.data;
        console.log(e);
    });*/

	$(function(){
        $("#txtTerminal").focus();
    });  	
    
</script>

<?php 

c::inicioFormulario("/ControlMigratorio/Aeropuertos/datosLector.php", "POST");
c::texto("txtTerminal","",135,200);
c::finFomulario();

echo "<br><br><br>";
c::finCentrado();
c::inicioCentrado();
    c::inicioTablaNormal($tablaAtributos);
    	c::inicioFila();
    		c::inicioColumna("align = left");
    			c::boton("Enviar", "Enviar al WebService", "ONCLICK=sendws()");
    		c::finColumna();
    	c::finFila();
        c::inicioFila();
            c::inicioColumna("align = left");
            	c::imagen("passport.png", "id='lamina' width='450x' height='300px' align='left' onload='conectar()'", "");
            c::finColumna();
            c::inicioColumna("align = right");
            	c::imagen("user.png", "id='face' width='200px' height='300px' align='right'", "");
            c::finColumna();
        c::finFila();
        c::inicioFila();
    		c::inicioColumna("align = left");
    			echo "Estado: "; c::texto("estado","",40,200,"readonly");
    		c::finColumna();
    		c::inicioColumna("align = right");
    			c::boton("Execute", "Capturar Foto","ONCLICK='camara()'");
    		c::finColumna();
    	c::finFila();
    c::finTabla();
c::finCentrado();
echo "<br><br><br>";

include 'DGME/estatico/pie.php';
?>
