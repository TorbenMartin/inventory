package main

//////////////////global content//////////////////
const style = `
<style>

#logo { position: absolute; top: 10px; left: 10px;}

button, input[type=button], input[type=submit], input[type=reset] {
  background-color: gray;
  border: none;
  color: white;
  text-decoration: none;
  cursor: pointer;
}

select[disabled=disabled], select:disabled,input[disabled=disabled], input:disabled {
  text-decoration: line-through;
  cursor: not-allowed;
  color: lightgray;
}

	
a:link {
color: gray;
}

a:visited {
color: gray;
}

a:hover {
color: 008000;
}

a:active {
color: gray;
}


.glass{
background: linear-gradient(135deg, rgba(255, 255, 255, 0.1), rgba(255, 255, 255, 0));
backdrop-filter: blur(10px);
-webkit-backdrop-filter: blur(10px);
border-radius: 10px;
border:1px solid rgba(255, 255, 255, 0.18);
box-shadow: 0 8px 32px 0 rgba(0, 0, 0, 0.37);
width: 700px;
min-height: 130px;
}

body {
margin: 0;
padding: 0;
overflow: scroll;
background-image: url("./img/bg1.jpg"); //add "" if you want
background-repeat: no-repeat;
background-attachment: fixed;
background-size: cover;
color: black;
}



.navbar {
  background-color: #034f84;
  height: 40px;
}

nav ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

nav {
  float: right;
  margin: 0;
}

li {
  display: inline;
  padding: 15px;
}

nav a {
  text-decoration: none;
  color: #353535;
  display: inline-block;
  padding-top: 15px;
}

nav a:hover {
  color: #00b2b2;
}

nav a:active {
  font-weight: bold;
}

@media screen and (max-width: 700px) {
  nav {
    width: 100%;
    margin-top: -5px;
  }

  nav li {
    display: block;
    background-color: #e5e5e5;
    text-align: center;
  }
 }

* {
  font-family: Montserrat, sans-serif;
}


summary {
   position: relative;  
   }

summary::marker {
   color: transparent;
}

summary::after {
   content:  "+"; 
   position: absolute;
   color: black;
   font-size: 3em;
   font-weight: bold; 
   right: 1em;
   top: .2em;
   transition: all 0.5s;
} 

details[open] summary::after {
 color: red;
 transform: translate(5px,0) rotate(45deg);
}


table#foobar-table > thead > tr {
 background-color: #e5e5e5;
 color: black;
 cursor: Pointer;
 font-size: 18px;
}




.obenkreis {
  width: 40px;
  height: 40px;
  border-radius: 50px;
  background-color: rgba(51,51,51,0.6);		
  animation: anitop 1s;
}
@keyframes anitop {
  0%{opacity:0}
	100%{opacity:1}
  }
.obenpfeil-1, .obenpfeil-2 {
  border: solid #fff;						
  border-width: 0 3px 3px 0;
  display: inline-block;
  padding: 5px;
    transform: rotate(-135deg);
    -webkit-transform: rotate(-135deg);
	position:absolute;
   left:50%;
   margin-left:-6px;
	transition: all 0.2s ease;
	}
.obenpfeil-1 {
	top:15px;
	}
.obenpfeil-2 {
	top:22px;
	}
.obenkreis:hover .obenpfeil-2 {
	top:8px;
}
#back-top2 {
position: fixed;
bottom: 5%;
right:5%;
z-index: 1000;
}
#back-top2  span{
display: block;
}
@media (max-width: 1680px) {
#back-top2 {
bottom: 5px;
right:5px;
}}
</style>
`

const globalmeta = `
<html lang="de">
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta http-equiv="Cache-Control" content="no-cache, no-store, must-revalidate" />
<meta http-equiv="Pragma" content="no-cache" />
<meta http-equiv="Expires" content="0" />
`

const menustart = `
<div class="navbar">
  <nav>
     
`

const menuend = `
    
  </nav>
</div>
`

var autologout string = "<meta http-equiv=\"refresh\" content=\"600; URL=/logout\">"

