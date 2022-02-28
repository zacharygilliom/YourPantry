let signUpForm = document.getElementById('sign-up-form');
let loginForm = document.getElementById('login-form');

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

async function signUpUser() {
	try {
		//username = document.getElementById('sign-up-email').value;	
		//password = document.getElementById('sign-up-password').value;
		let response = await fetch('http://localhost:8080/user/zacharygilliom@gmail.com/Penguin5');
		let data = await response.json();
	} catch (error) {
		console.log(error);
	}
}

loginForm.addEventListener('submit', function(event) {
	loginUser(event);
});
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
		console.log(data);
	} catch (error) {
		console.log(error);
	}
}
