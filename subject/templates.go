package subject

var (
	tmpIndex=`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>list</title>
    <style>
        #container{
            width: 50%;
            margin-left: 25%;
            margin-top: 15px;
        }
        table{border-right:1px solid #000;border-bottom:1px solid #000}
        table td{border-left:1px solid #000;border-top:1px solid #000}
        table th{border-left:1px solid #000;border-top:1px solid #000}
        a{
            text-decoration:none;
            color:black;
        }
        a:hover, a:link, a:active{color:blue;}
        table{
            width: 90%;
            background-color:#b0ff90;
            float:left;
        }
        .one{width:60px;}
        .three,.four{width:50px;}
        table th,table td{
            text-align: center;
        }
        #add{
            background: #b0ff90;
            width: 10%;
            float:left;
        }
        #layer{
            position: absolute ;
            background-color: black;
            z-index:1;
            top:0;
            bottom:0;
            left:0;
            right:0;
            filter:alpha(opacity=30);-moz-opacity:0.3;opacity: 0.3;
            display: none;
        }
        #editor{
            position: absolute ;
            width: 300px;
            height: 200px;
            background-color: #fff;
            z-index:2;
            left: 50%;
            margin-left: -150px;
            top:50%;
            margin-top: -150px;
            display: none;
        }
        #editor p{
            margin-left: 65px;
            margin-top: 25px;
        }
    </style>
    <script>
        var currentId,currentSubject;
        function test(){
            alert("test")
        }
        function showEditor(){
            document.getElementById("layer").style.display="inline";
            document.getElementById("editor").style.display="inline";
        }
        function hideEditor(){
            document.getElementById("layer").style.display="none";
            document.getElementById("editor").style.display="none";
        }
        function edit(id,subject){
            currentId=id;
            currentSubject=subject;
            document.getElementById("editor-id").value = id;
            document.getElementById("editor-subject").value = subject;
            showEditor();
        }
        function add(){
            currentId="";
            document.getElementById("editor-id").value = "";
            document.getElementById("editor-subject").value = "";
            showEditor();
        }
		function del(id){
            var xmlhttp=new XMLHttpRequest();
            xmlhttp.onreadystatechange=function()
            {
                if (xmlhttp.readyState==4 && xmlhttp.status==200)
                {
                    alert(xmlhttp.responseText);
                    if(xmlhttp.responseText == "success"){
                        location.reload();
                    }
                }
            }
            xmlhttp.open("POST","del",true);
            xmlhttp.setRequestHeader("Content-type","application/x-www-form-urlencoded");
            xmlhttp.send("id="+id);
        }
        function submit(){
            var id,subject;
            id = document.getElementById("editor-id").value;
            subject = document.getElementById("editor-subject").value;
            if(id == "" || subject==""){
                alert("id or subject is invalid!");
                return;
            }
            if(id == currentId&&subject==currentSubject){
                return;
            }
            var xmlhttp=new XMLHttpRequest();
            xmlhttp.onreadystatechange=function()
            {
                if (xmlhttp.readyState==4 && xmlhttp.status==200)
                {
                    alert(xmlhttp.responseText);
                    if(xmlhttp.responseText == "success"){
                        location.reload();
                    }
                }
            }
            xmlhttp.open("POST",currentId==""?"add":"update",true);
            xmlhttp.setRequestHeader("Content-type","application/x-www-form-urlencoded");
            xmlhttp.send("id="+id+"&subject="+subject+(currentId==""?"":"&oldId="+currentId));
        }
        function cancelEdit() {
            hideEditor();
        }
    </script>
</head>
<body>
<div id = "layer"></div>
<div id = "editor">
    <p>id<br><input id="editor-id" type="text" name="id" size="20" autocomplete="off"></p>
    <p>subject<br><input id= "editor-subject" type="text" name="subject" size="20" autocomplete="off"></p>
    <p><button onclick="submit()">submit</button>&nbsp;<button onclick="cancelEdit()">cancel</button></p>
</div>
<div id="container">
    <table id="list" border="0" cellspacing="0" cellpadding="0">
        <tr>
            <th class="one">id</th>
            <th class="two">subject</th>
            <th class="three"></th>
            <th class="four"></th>
        </tr>
        {{ range $key, $value := . }}
			<tr>
			<td>{{ $key }}</td>
			<td>{{ $value }}</td>
			<td><a href="javascript:edit({{ $key }},'{{ $value }}')">edit</a></td>
			<td><a href="javascript:del({{ $key }})">delete</a></td>
			</tr>
		{{ end }}
    </table>
    <button id="add" onclick="add()">add</button>
</div>
</body>
</html>`


	tmpLogin=`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>login</title>
</head>
<body>
<form method="post" action="loginservice">
  <p>name<br><input type="text" name="name" size="20"></p>
  <p>passwold<br><input type="password" name="password" size="20"></p>
  <p><input type="submit" value="login">
  </form>
</body>
</html>`
)