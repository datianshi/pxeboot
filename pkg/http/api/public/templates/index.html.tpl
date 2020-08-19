<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <meta name="description" content="">
    <meta name="author" content="Mark Otto, Jacob Thornton, and Bootstrap contributors">
    <meta name="generator" content="Jekyll v4.1.1">
    <title>Starter Template · Bootstrap</title>

    <link rel="canonical" href="https://getbootstrap.com/docs/4.5/examples/starter-template/">

    <!-- Bootstrap core CSS -->
    <link href="static/assets/bootstrap-4.5.2/dist/css/bootstrap.min.css" rel="stylesheet">

    <style>
        .bd-placeholder-img {
            font-size: 1.125rem;
            text-anchor: middle;
            -webkit-user-select: none;
            -moz-user-select: none;
            -ms-user-select: none;
            user-select: none;
        }

        @media (min-width: 768px) {
            .bd-placeholder-img-lg {
                font-size: 3.5rem;
            }
        }
    </style>
    <!-- Custom styles for this template -->
    <link href="static/assets/css/starter-template.css" rel="stylesheet">
    <script src="https://code.jquery.com/jquery-3.5.1.min.js"></script>
    <script src="static/assets/bootstrap-4.5.2/dist/js/bootstrap.bundle.min.js"></script>
    <script src="static/assets/js/start-template.js"></script>
</head>
<body>
<nav class="navbar navbar-expand-md navbar-dark bg-dark fixed-top">
</nav>

<main role="main" class="container">

    <div class="row">
        <div class="col-md">
            <form id="mac_edit_form">
                <div class="form-group" >
                    <label for="mac_address_edit">Mac Address</label>
                    <select class="form-control" id="mac_address_edit">
                    </select>
                </div>
                <div class="form-group">
                    <label for="static_ip_edit">Static IP</label>
                    <input class="form-control" id="static_ip_edit" placeholder="10.65.62.100">
                </div>
                <div class="form-group">
                    <label for="host_name_edit">Host Name</label>
                    <input class="form-control" id="host_name_edit" placeholder="server.example.com">
                </div>
                <button id="edit_nic_button" class="btn btn-primary">Submit</button>
                <button id="delete_nic_button" class="btn btn-primary">Delete</button>
            </form>
        </div>
        <div class="col-md">
            <form id="mac_add_form">
                <div class="form-group">
                    <label for="mac_address_add">Mac Address</label>
                    <input class="form-control" id="mac_address_add" placeholder="example: 00:50:A6:83:75:98">
                    </select>
                </div>
                <div class="form-group">
                    <label for="static_ip_add">Static IP</label>
                    <input class="form-control" id="static_ip_add" placeholder="example: 10.65.34.67">
                </div>
                <div class="form-group">
                    <label for="host_name_add">Host Name</label>
                    <input class="form-control" id="host_name_add" placeholder="example: server1.example.org">
                </div>
                <button type="submit" class="btn btn-primary">Add</button>
            </form>
        </div>
    </div>

    <div class="row config_content">
        <form id="upload_image" action="api/image" method="post" enctype="multipart/form-data">
            <div class="form-group">
                <input type="file" class="form-control" id="upload_image_input" accept=".iso" name="image">
                <button type="submit" class="btn btn-primary">Upload</button>
            </div>
        </form>
        <div>

    <div class="row config_content">
        <pre>
            <code id="config_content">
            </code>
        </pre>
    </div>

</main><!-- /.container -->
</html>
