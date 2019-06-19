var http = require('http');

http.createServer( (request, response) => {
    response.writeHead(200);
    setTimeout(()=>{
        response.end(new Array(1024).fill('haha').join(''));
	}, 50);
}).listen(3000);
