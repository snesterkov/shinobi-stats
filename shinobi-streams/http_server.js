 var connect = require('connect');
 var serveStatic = require('serve-static');

 connect()
     .use(serveStatic('/home/www'))
     .listen(8000, () => console.log('Server running on 8000...'));