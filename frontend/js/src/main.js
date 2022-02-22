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
		username = document.getElementById('sign-up-email').value;	
		password = document.getElementById('sign-up-password').value;
		let response = await fetch('http://localhost:8080/username/password');
		let data = await response.json();
	} catch (error) {
		console.log(error);
	}
}

async function loginUser() {
	try {
		username = document.getElementById('login-email').value;
		password = document.getElementById('login-password').value;
		var apiUrl = "http://localhost:8080/login/" + username + "/" + password
		console.log(apiUrl);
		let response = await fetch(apiUrl);
		console.log(response);
		let data = await response.json();
		console.log(data)
	} catch (error) {
		console.log(error);
	}
}

function main() {
	//fetchIngredList()
}

main()
