async function fetchIngredList() {
	try {
		let response =  await fetch('http://localhost:8080/61ece6d2e84c62bdcdbcc42d/ingredients/list');
		let data = await response.json();
		console.log(data["ingredients"])
		for (let i = 0; i < data['ingredients'].length; i++) {
			console.log(data['ingredients'][i])
		}
	} catch (error) {
		console.log(error);
	}
}

fetchIngredList()
