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
/*
async function loginUser() {
	try {
		//username = document.getElementById('login-email').value;
		//password = document.getElementById('login-password').value;
		//console.log(username);
		//console.log(password);
		let response = await fetch('http://localhost:8080/login?email=zacharygilliom@gmail.com&password=Penguin5');
		console.log(response);
		let data = await response.json();
		console.log(data);
	} catch (error) {
		console.log(error);
	}
}
*/
function loginUser() {
	fetch("http://localhost:8080/login?email=zacharygilliom@gmail.com&password=Penguin5")
	.then(response => {
		console.log(response);
	})
	.catch(error => {
		console.log(error);
	})
}
async function autoLoginUser() {
	try {
		let response = await fetch('http://localhost:8080/login?email=zacharygilliom@gmail.com&password=Penguin5');
		console.log(response);
		let data = await response.json();
		console.log(data);
	} catch (error) {
		console.log(error);
	}
}


//autoLoginUser()

