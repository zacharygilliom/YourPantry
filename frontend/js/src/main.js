function ingredients() {
	fetchIngredList();
	let quickAddForm = document.getElementById('quick-add-form');
	quickAddForm.addEventListener('submit', function(event) {
		quickAddIngredient(event);
	});
	
}
function landing() {
	let signUpForm = document.getElementById('sign-up-form');
	let loginForm = document.getElementById('login-form');
	loginForm.addEventListener('submit', function(event) {
		loginUser(event);
	});
	signUpForm.addEventListener('submit', function(event) {
		signUpUser(event);
	});
}

async function fetchIngredList() {
	try {
		var token = getCookie("token");
		const requestOption = {
			method:'GET',
			headers: {'Content-Type': 'application/json', 'Authorization':'Bearer ' + token},
			SendCookie: true,
			SecureCookie: false,
			CookieDomain: "localhost:8080",
			CookieName: "token",
			TokenLookup: "cookie:token",
			credentials:'include'
		};
		let response = await fetch('http://localhost:8080/user/ingredients/list', requestOption);
		let data = await response.json();
		if (data['size'] > 0) {
			ul = document.createElement('ul');
			ul.className = 'list-group';
			ul.id = 'ingredient-list-items';
			document.getElementById('ingredient-list').appendChild(ul);
			let idx = 0;
			data['ingredients'].forEach(function (item){
				createIngredientList(ul, item, idx);
				idx += 1;
			})
		} else if (data['size'] == 0) {
			item = document.getElementById('ingredient-list');
			p = document.createElement('p');
			p.innerHTML = "Nothing to show!";
			item.appendChild(p);
		}
	} catch (error) {
		console.log(error);
	}
}

function createIngredientList(ul, item, idx) {
	let li  = document.createElement('li');
	li.className ='list-group-item d-flex justify-content-between align-items-center';
	li.id = 'list-group-item ' + idx;
	ul.appendChild(li);
	li.innerHTML += item;
	let bt = document.createElement('button');
	bt.className = "btn btn-outline-danger";
	bt.type = "button";
	bt.id = 'remove-ingredient-button ' + idx;
	bt.onclick =  function(event) {
		removeIngredient(event);
	};
	li.appendChild(bt);
	let iconSVG = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
	let path1SVG = document.createElementNS('http://www.w3.org/2000/svg', 'path');
	let path2SVG = document.createElementNS('http://www.w3.org/2000/svg', 'path');
	iconSVG.setAttribute('fill', 'currentColor');
	iconSVG.setAttribute('viewBox', '0 0 16 16');
	iconSVG.setAttribute('width', '16');
	iconSVG.setAttribute('height', '16');
	path1SVG.setAttribute('d',"M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0V6z");
	path2SVG.setAttribute('fill-rule', 'evenodd');
	path2SVG.setAttribute('d', "M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1v1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4H4.118zM2.5 3V2h11v1h-11z");
	bt.appendChild(iconSVG);
	iconSVG.appendChild(path1SVG);
	iconSVG.appendChild(path2SVG);
	//return li
}

async function signUpUser(event) {
	try {
		event.preventDefault();
		username = document.getElementById('sign-up-email').value;	
		pass = document.getElementById('sign-up-password').value;
		fname = document.getElementById('sign-up-fname').value;
		lname = document.getElementById('sign-up-lname').value;
		let userData = {email:username, password:pass, firstname: fname, lastname:lname};
		//console.log(userData);
		const requestOption = {
			method: 'POST',
			headers: {'Content-Type': 'application/json'},
			body: JSON.stringify(userData)
	};
		let response = await fetch('http://localhost:8080/sign-up', requestOption);
		let data = await response.json();
		//console.log(data);
		if (data["data"] == 1) {
			window.location.replace("landing.html");
			alert("New User Account has been created.  Please sign in Below!")
			return false
		}
	} catch (error) {
		console.log(error);
	}
}

async function loginUser(event) {
	try {
		event.preventDefault();
		username = document.getElementById('login-email').value;
		pass = document.getElementById('login-password').value;
		let userData = {email:username, password:pass};
		const requestOption = {
			method: 'POST',
			headers: {'Content-Type': 'application/json'},
			body: JSON.stringify(userData)
		};
		let response = await fetch('http://localhost:8080/login', requestOption);
		let data = await response.json();
		if (data["code"] != 200) {
			alert("Login Not Successful: Please Try Again");
			location.reload();
		} else if (data["code"] == 200) {
			document.cookie = "token=" + data['token'] + "; path=/; SameSite=None; secure=true;"
			window.location.replace("home.html");
			return false;
		}
	//TODO: Catch the Error
	} catch (error) {
		console.log(error);
	}
}

async function quickAddIngredient(event) {
	try {
		event.preventDefault();
		ingredient = document.getElementById('Ingredient-selection').value;
		let userData = {ingredient:ingredient};
		var token = getCookie("token");
		const requestOption = {
			method:'POST',
			headers: {'Content-Type': 'application/json', 'Authorization':'Bearer ' + token},
			body: JSON.stringify(userData),
			SendCookie: true,
			SecureCookie: false,
			CookieDomain: "localhost:8080",
			CookieName: "token",
			TokenLookup: "cookie:token",
			credentials:'include'
		};
		let response = await fetch('http://localhost:8080/user/ingredients/add', requestOption);
		let data = await response.json();
		if (data['code'] != 200) {
			//alert("Ingredient not added! Please Try Again");
			location.reload();
		} else if (data['code'] == 200) {
			//alert("Ingredient has been added!");
			location.reload();
		}
	} catch (error) {
		console.log(error);
	}
}

async function removeIngredient(event) {
	try {
		let bt_clicked_id = String(event.currentTarget.id).charAt(String(event.currentTarget.id).length-1);
		ingredient = document.getElementById('list-group-item ' + bt_clicked_id).textContent;
		let userData = {ingredient:ingredient};
		var token = getCookie("token");
		const requestOption = {
			method:'POST',
			headers: {'Content-Type': 'application/json', 'Authorization': 'Beader ' + token},
			body: JSON.stringify(userData),
			SendCookie: true,
			SecureCookie: false,
			CookieDomain: 'localhost:8080',
			CookieName: 'token',
			TokenLookup: 'cookie:token',
			credentials: 'include'
		};
		let response = await fetch('http://localhost:8080/user/ingredients/remove', requestOption);
		let data = await response.json();
		if (data['code'] == 200) {
			location.reload();
		}
	} catch (error) {
		location.reload();
		console.log(error);
	}
}

function getCookie(cname) {
	let name = cname + "=";
	let decodedCookie = decodeURIComponent(document.cookie);
	let ca = decodedCookie.split(';');
	for (let i = 0; i < ca.length; i++) {
		let c = ca[i];
		while (c.charAt(0) == ' ') {
			c = c.substring(1);
		}
		if (c.indexOf(name) == 0) {
			return c.substring(name.length, c.length);
		}
	}
	return "";
}
