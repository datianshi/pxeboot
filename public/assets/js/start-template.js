// Material Select Initialization
$(document).ready(function() {
    $("button#submit").click(function(){
        $.ajax({
            type: "GET",
            url: "PastSurgicalCustomItem",
            data: $('form.form-horizontal').serialize(),
            success: function(msg){
                alert(msg);
            },
            error: function(){
                alert("failure");
            }
        });
    });
    $("#exampleFormControlSelect1").click(function () {
        console.log("hello")
        $(this).append("<option>abc</option>")
    });
})