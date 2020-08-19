// Material Select Initialization
$(document).ready(function() {
    load_nics()
    select_nic_change()
    submit_nic_change()
    add_nic_change()
    delete_nic()
    load_conf()
})

function load_nics() {
    $.get("api/conf/nics", function(data, status, xhr){
        $("#mac_address_edit").empty()
        if (data.length > 0) {
            for (let i = 0; i < data.length; i++) {
                $("#mac_address_edit").append("<option>" + data[i].mac_address + "</option>")
            }
            load_nic(data[0])
        }
    });
}

function load_conf() {
    $.get("api/conf", function(data, status){
        $("#config_content").text(JSON.stringify(data, undefined, 4))
    });
}

function load_nic(nic) {
    $("#host_name_edit").val(nic.hostname)
    $("#static_ip_edit").val(nic.ip)
    $("#mac_address_edit").val(nic.mac_address)
}

function select_nic_change(){
    $("#mac_address_edit").change(function () {
        var current_mac = $(this).val()
        $.get("api/conf/nic/" + current_mac, function(data, status){
            load_nic(data)
        })
    })
}

function submit_nic_change(){
    $("#edit_nic_button").click(function() {
        var putData = {}
        putData.hostname = $("#host_name_edit").val()
        putData.ip = $("#static_ip_edit").val()
        var current_mac = $("#mac_address_edit").val()
        $.ajax({
            url: "api/conf/nic/" + current_mac,
            method: 'PUT',
            data: JSON.stringify(putData),
            contentType: 'application/json;charset=UTF-8',
            success: function(data) {
                alert("updated" + current_mac)
            }
        })

    })
}

function delete_nic(){
    $("#delete_nic_button").click(function (){
        var current_mac = $("#mac_address_edit").val()
        $.ajax({
            url: "api/conf/nic/" + current_mac,
            method: 'DELETE',
            contentType: 'application/json;charset=UTF-8',
            success: function(data) {
                alert("delete " + current_mac)
            }
        })

    })
}

function add_nic_change(){
    $("#mac_add_form").submit(function (e) {
        e.preventDefault()
        var postData = {}
        postData.hostname = $("#host_name_add").val()
        postData.ip = $("#static_ip_add").val()
        postData.mac_address = $("#mac_address_add").val()
        $.ajax({
            url: "api/conf/nic",
            method: 'POST',
            data: JSON.stringify(postData),
            contentType: 'application/json;charset=UTF-8',
            success: function(data) {
                alert("new nic configured")
                load_nics()
                load_conf()
            },
            error: function(xhr, status, error){
                alert('Error - ' + xhr.responseText);
            }
        })

    })
}