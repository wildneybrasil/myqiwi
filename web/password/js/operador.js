function setupOperator(operatorId, targetFormId) {
    $.ajax({
        url: "/data/operator",
        async: true,
        data: { operatorId: operatorId },
        type: "POST",
        crossDomain: false,
        success: function(data) {
            var jsonData = JSON.parse( data );

            $("#" + targetFormId + " [name='name']").val(jsonData.data.name);
            $("#" + targetFormId + " [name='email']").val(jsonData.data.email);
            $("#" + targetFormId + " [name='operatorId']").val(jsonData.data.id);

            $("#" + targetFormId).siblings("form").hide();
            $("#" + targetFormId).fadeIn();
            console.log($("#" + targetFormId));
        },
        error: function(err) {
            console.log("Error while accessing: " + table.url);
            console.log(err);
        }
    });
}

function updateOperator() {
    alert( $('#name').val() );

    $.ajax({
        url: "/update/operator",
        async: true,
        data: { operatorId: $('#operatorId').val(),  name:$('#name').val(), password: $('#password').val(), email: $('#email').val()   },
        type: "POST",
        crossDomain: false,
        success: function(data) {
            alert("Dados alterados com sucesso")
        },
        error: function(err) {
            console.log("Error while accessing: " + table.url);
            console.log(err);
        }
    });
}
