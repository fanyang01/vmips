.text
	main:
	li $v0, 4
	la $a0, hello
	syscall
	li $v0, 10
	syscall

.data
	hello:
	.ascii "Hello, world!"
	.byte 0x0A, 0
