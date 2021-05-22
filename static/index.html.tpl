<html>

<head>
  <title>k8s version</title>
  <meta charset="UTF-8">
  <link rel="stylesheet" type="text/css" href="style.css"/>
  <link rel="preconnect" href="https://fonts.gstatic.com">
  <link href="https://fonts.googleapis.com/css2?family=Roboto:wght@500&display=swap" rel="stylesheet">
  <script src="https://kit.fontawesome.com/9132f35dd2.js" crossorigin="anonymous"></script>
</head>

<body>
  <div class="grid">
    {{ range $i, $e := .VersionSources }}
    <div class="grid-item" style="background-color:{{.Color}};{{if eq $i 0}}width:100%;{{end}}">
      <p>{{if eq $i 1}}<i class="fas fa-medal"></i>{{end}} {{.Name}}</p>
      <h1>{{.Version}}</h1>
    </div>
    {{ end }}
  </div>
  <div class="footer">
    <i class="fas fa-home"></i> Lainecloud | <i class="fab fa-github"></i> Github
  </div>
</body>

</html>
