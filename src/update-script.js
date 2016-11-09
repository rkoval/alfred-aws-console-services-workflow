// this is a hacked-together script that will pluck out service links and their descriptions to JSON
// for later adding to the csv
//
// to use it, run it on the AWS homepage dashboard
var texts = $('.servicesContainer').find('a').find('.serviceName, .serviceCaption').map((a, b) => {
  return {
    text: b.textContent,
    href: b.parentElement.href
  }
});
var outputs = []
var string;
for (var i = 0; i < texts.length; i++) {
  var obj = texts[i];
  if (i % 2 === 0) {
    string = obj.text;
  } else {
    var {text, href} = obj;
    var regex = /\.com\/(.*?)\//g;
    var urlKey = regex.exec(href)[1];
    var formattedHref = `https://console.aws.amazon.com/${urlKey}/home`
    var output = `"${urlKey}","${string} - ${text}","${formattedHref}"`;
    outputs.push(output);
  }
}
JSON.stringify(outputs)
