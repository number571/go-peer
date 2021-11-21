function copy_text(id) {
	var copyText = document.getElementById(id);
	copyText.select();
	document.execCommand("copy");
	open_block('message_copy');
}

function view_block(id) {
	var e = document.getElementById(id);
	if(e.style.display == 'block')
		e.style.display = 'none';
	else
		e.style.display = 'block';
}

function open_block(id) {
	var e = document.getElementById(id);
	e.style.display = 'block';
}

function close_block(id) {
	var e = document.getElementById(id);
	e.style.display = 'none';
}

function clear_value(id) {
    var e = document.getElementById(id);
    e.value = "";
}
