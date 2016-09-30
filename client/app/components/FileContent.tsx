import * as React from 'react';


export function FileContent(props) {
    const code = `
$('#loading-example-btn').click(function () {
var btn = $(this)
btn.button('loading')
$.ajax(...).always(function () {
    btn.button('reset')
});
});
    `
    return (
        <pre data-enlighter-language='js'>
            {code}
        </pre>
    );
}