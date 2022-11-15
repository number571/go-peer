function copy_text(id) {
	var copyText = document.getElementById(id);
	copyText.select();
	document.execCommand("copy");
}

function open_block(id) {
	var e = document.getElementById(id);
	e.style.display = 'block';
}

function close_block(id) {
	var e = document.getElementById(id);
	e.style.display = 'none';
}
