var http = require('http');

http.createServer( (request, response) => {
    response.writeHead(200);
    setTimeout(()=>{
        response.end(new Array(4096).fill('s').join(''));
	}, 50);
}).listen(8888);
