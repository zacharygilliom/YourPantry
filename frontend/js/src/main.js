function ingredients() {
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
		let response =  await fetch('http://localhost:8080/61ece6d2e84c62bdcdbcc42d/ingredients/list');
		let data = await response.json();
		/*let data = await {'ingredients': [],
		}
		*/
		var str = '<ol>'
		data['ingredients'].forEach(function(ingredient) {
			str += '<li>' + ingredient + '</li>'
		});
		str += '</ol>';
		document.getElementById('ingredient-list').innerHTML = str;
	} catch (error) {
		console.log(error);
	}
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
		console.log(response);
		console.log(data);
		if (data["code"] != 200) {
			alert("Login Not Successful: Please Try Again");
			//location.reload();
		} else if (data["code"] == 200) {
			document.cookie = "token=" + data['token'] + "; path=/; SameSite=None; secure=true;"
			//window.localStorage.setItem('token', data['token']);
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
		console.log(token);
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
		//console.log(userData);
		let response = await fetch('http://localhost:8080/user/ingredients/add', requestOption);
		let data = await response.json();
		console.log(response);
		console.log(data);
	} catch (error) {
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
