$(document).ready(function(){

	var user_name
	var final_conexion
	$("#form_registro").on("submit",function(e){
		console.log("Entro")
		e.preventDefault();
		user_name = $("#user_name").val()

		$.ajax({
			type:"POST",
			url:"https://chatmultiusos.herokuapp.com/validate",
			data:{
				"user_name": user_name
			},

			success: function(data){
				result(data)
			}

		})
	})

	function result(data){
		obj = JSON.parse(data)
		if(obj.isvalid === true){
			create_conexion()

		}else{
			console.log("Intentalo de nuevo")
		}
	}

	function create_conexion(){
		$("#registro").hide()
		$("#container_chat").show()

		var conexion = new WebSocket("wss://chatmultiusos.herokuapp.com/chat/"+ user_name)
		final_conexion = conexion
        conexion.onopen = function(response){
        	conexion.onmessage = function(response){
        		console.log(response.data)
        		val = $("#chat_area").val()
        		$("#chat_area").val(val + "\n"+response.data)
        	}
        }
	}

	$("#form_message").on("submit", function(e){
		e.preventDefault();
		mensaje = $("#msg").val()
		final_conexion.send(mensaje)
		$("#msg").val("")
	})
})