// Material Select Initialization
$(document).ready(function() {
    load_nics()
    select_nic_change()
})

function load_nics() {
    $("#mac_address").load("api/conf/nics", function(data, status, xhr){
        var nics = JSON.parse(data)
        if (nics.length > 0) {
            for (let i = 0; i < nics.length; i++) {
                $(this).append("<option>" + nics[i].mac_address + "</option>")
            }
            load_nic(nics[0])
        }
    });
}

function load_nic(nic) {
    $("#dhcp_ip").val(nic.dhcp_ip)
    $("#host_name").val(nic.hostname)
    $("#static_ip").val(nic.ip)
}

function select_nic_change(){
    $("#mac_address").change(function () {
        var current_mac = $(this).val()
        $.get("api/conf/nic/" + current_mac, function(data, status){
            load_nic(data)
        })
    })
}