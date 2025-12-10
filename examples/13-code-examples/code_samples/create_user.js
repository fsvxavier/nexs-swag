// JavaScript example for creating a user
const createUser = async () => {
  const response = await fetch('http://localhost:8080/api/users', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      name: 'John Doe'
    })
  });
  
  const data = await response.json();
  console.log(data);
};

createUser();
