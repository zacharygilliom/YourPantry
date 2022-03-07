let signUpForm = document.getElementById('sign-up-form');
let loginForm = document.getElementById('login-form');

loginForm.addEventListener('submit', function(event) {
	loginUser(event);
});

signUpForm.addEventListener('submit', function(event) {
	signUpUser(event);
});

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
			window.location.replace("home.html");
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
		//console.log(data);
		if (data["data"] == 0) {
			alert("Login Not Successful: Please Try Again");
			//location.reload();
		} else if (data["data"] == 1) {
			window.location.replace("home.html");
			return false;
		}
	} catch (error) {
		console.log(error);
	}
}
