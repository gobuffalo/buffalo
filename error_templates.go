package buffalo

var devErrorTmpl = `
<html>
<head>
  <title><%= status %> - ERROR!</title>
  <style>html{font-family:sans-serif;-webkit-text-size-adjust:100%;-ms-text-size-adjust:100%}body{margin:0}header{display:block}a{background-color:transparent}a:active,a:hover{outline:0}h1{margin:.67em 0;font-size:2em}img{border:0}pre{overflow:auto}code,pre{font-family:monospace,monospace;font-size:1em}table{border-spacing:0;border-collapse:collapse}td,th{padding:0}@media print{*{color:#000!important;text-shadow:none!important;background:0 0!important;-webkit-box-shadow:none!important;box-shadow:none!important}a,a:visited{text-decoration:underline}a[href]:after{content:" (" attr(href) ")"}pre{border:1px solid #999;page-break-inside:avoid}thead{display:table-header-group}img,tr{page-break-inside:avoid}img{max-width:100%!important}h3{orphans:3;widows:3}h3{page-break-after:avoid}.table{border-collapse:collapse!important}.table td,.table th{background-color:#fff!important}}@font-face{font-family:'Glyphicons Halflings';src:url(../fonts/glyphicons-halflings-regular.eot);src:url(../fonts/glyphicons-halflings-regular.eot?#iefix) format('embedded-opentype'),url(../fonts/glyphicons-halflings-regular.woff2) format('woff2'),url(../fonts/glyphicons-halflings-regular.woff) format('woff'),url(../fonts/glyphicons-halflings-regular.ttf) format('truetype'),url(../fonts/glyphicons-halflings-regular.svg#glyphicons_halflingsregular) format('svg')}*{-webkit-box-sizing:border-box;-moz-box-sizing:border-box;box-sizing:border-box}html{font-size:10px;-webkit-tap-highlight-color:rgba(0,0,0,0)}body{font-family:"Helvetica Neue",Helvetica,Arial,sans-serif;font-size:14px;line-height:1.42857143;color:#333;background-color:#fff}a{color:#337ab7;text-decoration:none}a:focus,a:hover{color:#23527c;text-decoration:underline}a:focus{outline:5px auto -webkit-focus-ring-color;outline-offset:-2px}img{vertical-align:middle}h1,h3{font-family:inherit;font-weight:500;line-height:1.1;color:inherit}h1,h3{margin-top:20px;margin-bottom:10px}h1{font-size:36px}h3{font-size:24px}code,pre{font-family:Menlo,Monaco,Consolas,"Courier New",monospace}code{padding:2px 4px;font-size:90%;color:#c7254e;background-color:#f9f2f4;border-radius:4px}pre{display:block;padding:9.5px;margin:0 0 10px;font-size:13px;line-height:1.42857143;color:#333;word-break:break-all;word-wrap:break-word;background-color:#f5f5f5;border:1px solid #ccc;border-radius:4px}.container{padding-right:15px;padding-left:15px;margin-right:auto;margin-left:auto}@media (min-width:768px){.container{width:750px}}@media (min-width:992px){.container{width:970px}}@media (min-width:1200px){.container{width:1170px}}.row{margin-right:-15px;margin-left:-15px}.col-md-1,.col-md-10,.col-md-12,.col-sm-2,.col-sm-6,.col-xs-3,.col-xs-7{position:relative;min-height:1px;padding-right:15px;padding-left:15px}.col-xs-3,.col-xs-7{float:left}.col-xs-7{width:58.33333333%}.col-xs-3{width:25%}@media (min-width:768px){.col-sm-2,.col-sm-6{float:left}.col-sm-6{width:50%}.col-sm-2{width:16.66666667%}}@media (min-width:992px){.col-md-1,.col-md-10,.col-md-12{float:left}.col-md-12{width:100%}.col-md-10{width:83.33333333%}.col-md-1{width:8.33333333%}}table{background-color:transparent}th{text-align:left}.table{width:100%;max-width:100%;margin-bottom:20px}.table>tbody>tr>td,.table>thead>tr>th{padding:8px;line-height:1.42857143;vertical-align:top;border-top:1px solid #ddd}.table>thead>tr>th{vertical-align:bottom;border-bottom:2px solid #ddd}.table>thead:first-child>tr:first-child>th{border-top:0}.table-striped>tbody>tr:nth-of-type(odd){background-color:#f9f9f9}.container:after,.container:before,.row:after,.row:before{display:table;content:" "}.container:after,.row:after{clear:both}@-ms-viewport{width:device-width}
	h1{margin-top:20px}*{-webkit-box-sizing:border-box;-moz-box-sizing:border-box;box-sizing:border-box}body{font-family:"Helvetica Neue",Helvetica,Arial,sans-serif;font-size:14px;line-height:1.42857143;color:#333;background-color:#fff;margin:0}h1{margin-bottom:10px;font-family:inherit;font-weight:500;line-height:1.1;color:inherit}.table{margin-bottom:20px}h1{font-size:36px}a{color:#337ab7;text-decoration:none}a:hover{color:#23527c}.container{padding-right:15px;padding-left:15px;margin-right:auto;margin-left:auto}@media (min-width:768px){.container{width:750px}}@media (min-width:992px){.container{width:970px}}@media (min-width:1200px){.container{width:1170px}}.table{width:100%;max-width:100%;background-color:transparent;border-spacing:0;border-collapse:collapse}.table-striped>tbody{background-color:#f9f9f9}.table>tbody>tr>td,.table>thead>tr>th{padding:8px;line-height:1.42857143;vertical-align:top;border-top:1px solid #ddd}.table>thead>tr>th{border-top:0;vertical-align:bottom;border-bottom:2px solid #ddd;text-align:left}code{padding:2px 4px;font-size:90%;color:#c7254e;background-color:#f9f2f4;border-radius:4px;font-family:Menlo,Monaco,Consolas,"Courier New",monospace}.row{margin-right:-15px;margin-left:-15px}.col-md-10{float:left;position:relative;min-height:1px;padding-right:15px;padding-left:15px}.col-md-10{width:83.33333333%}img{vertical-align:middle;border:0}.container{min-width:320px}body{font-family:helvetica}table{font-size:14px}table.table tbody tr td{border-top:0;padding:10px}pre{white-space:pre-line;margin-bottom:10px;max-height:275px;overflow:scroll}header{background-color:#ed605e;padding:10px 20px;box-sizing:border-box}.logo img{width:80px}.titles h1{font-size:30px;font-weight:300;color:#fff;margin:24px 0}.content h3{color:gray;margin:25px 0}.foot{padding:5px 0 20px;text-align:right;color:#c5c5c5;font-weight:300}.foot a{color:#8b8b8b;text-decoration:underline}.centered{text-align:center}@media all and (max-width:500px){.titles h1{font-size:25px;margin:26px 0}}@media all and (max-width:530px){.titles h1{font-size:20px;margin:24px 0}.logo{padding:0}.logo img{width:100%;max-width:80px}}
  </style>
</head>

<body>
  <header>
    <div class="container">
      <div class="row">
        <div class="col-md-1 col-sm-2 col-xs-3 logo">
          <a href="/"><img src="https://gobuffalo.io/assets/images/logo_med.png" alt=""></a>
        </div>
        <div class="col-md-10 col-sm-6 col-xs-7 titles">
          <h1>
            <%= status %> - ERROR!
          </h1>
        </div>
      </div>
    </div>
  </header>

  <div class="container content">
    <div class="row">
      <div class="col-md-12">
        <h3>Error Trace</h3>
        <pre><%= error %></pre>

        <h3>Context</h3>
        <pre><%= inspect(context) %></pre>

        <h3>Parameters</h3>
        <pre><%= inspect(params) %></pre>

        <h3>Headers</h3>
        <pre><%= inspect(headers) %></pre>

        <h3>Form</h3>
        <pre><%= inspect(posted_form) %></pre>

        <h3>Routes</h3>
        <table class="table table-striped">
          <thead>
            <tr text-align="left">
              <th class="centered">METHOD</th>
              <th>PATH</th>
              <th>NAME</th>
              <th>HANDLER</th>
            </tr>
          </thead>
          <tbody>

            <%= for (r) in routes { %>
              <tr>
                <td class="centered">
                  <%= r.Method %>
                </td>
                <td>
                  <%= if (r.Method != "GET" || r.Path ~= "{") { %>
                    <%= r.Path %>
                  <% } else { %>
                    <a href="<%= r.Path %>"><%= r.Path %></a>
                  <% } %>
                </td>
                <td>
                  <%= r.PathName %>
                </td>
                <td><code><%= r.HandlerName %></code></td>
              </tr>
            <% } %>

          </tbody>
        </table>
      </div>
    </div>
    <div class="foot"> <span> Powered by <a href="http://gobuffalo.io/">gobuffalo.io</a></span></div>
  </div>
</body>
</html>
`
var prodErrorTmpl = `
<!DOCTYPE html>
<html>
<head>
<style>h1,p.powered{text-align:center}body{background:#ECECEC;padding-top:25px;font-family:helvetica neue,helvetica,sans-serif;color:#333}.card{box-sizing:border-box;width:440px;min-width:270px;margin:0 auto;padding:10px 25px 35px 10px;background:#FFF;box-shadow:0 2px 4px 0 rgba(185,185,185,.28);border-radius:5px}.card p{max-width:320px;margin:15px auto}h1{font-size:22px}hr{border:.5px solid #D72727;width:180px}p.powered{font-family:HelveticaNeue-Light;font-size:12px;color:#333}@media (max-width:600px){.card{width:100%;display:block}}</style>
</head>
<body>
<div class="container">
	<div class="card">
		<h1>We're Sorry!</h1>
		<hr>
		<p>It looks like something went wrong! Don't worry, we are aware of the problem and are looking into it.</p>
		<p>Sorry if this has caused you any problems. Please check back again later.</p>
	</div>

	<p class="powered">powered by <a href="https://gobuffalo.io">gobuffalo.io</a></p>
</div>
</body>
</html>
`

var prodNotFoundTmpl = `
<!DOCTYPE html>
<html>
<head>
<style>h1,p.powered{text-align:center}body{background:#ECECEC;padding-top:25px;font-family:helvetica neue,helvetica,sans-serif;color:#333}.card{box-sizing:border-box;width:440px;min-width:270px;margin:0 auto;padding:10px 25px 35px 10px;background:#FFF;box-shadow:0 2px 4px 0 rgba(185,185,185,.28);border-radius:5px}.card p{max-width:320px;margin:15px auto}h1{font-size:22px}hr{border:.5px solid #1272E2;width:180px}p.powered{font-family:HelveticaNeue-Light;font-size:12px;color:#333}@media (max-width:600px){.card{width:100%;display:block}}</style>
</head>
<body>
<div class="container">
	<div class="card">
		<h1>Not Found</h1>
		<hr>
		<p>The page you're looking for does not exist, you may have mistyped the address or the page may have been moved.</p>
	</div>

	<p class="powered">powered by <a href="https://gobuffalo.io">gobuffalo.io</a></p>
</div>
</body>
</html>
`
