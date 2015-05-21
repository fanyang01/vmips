.text
.globl main
main:
	li $v0, 4
	la $a0, prompt
	syscall

	li $v0, 5
	syscall

# print "fib(n) = "
	move $s0, $v0
	li $v0, 4
	la $a0, result_a
	syscall
	move $a0, $s0
	li $v0, 1
	syscall
	la $a0, result_b
	li $v0, 4
	syscall

	move $a0, $s0 # call function fib
	jal fib
	move $a0, $v0 # print result
	li $v0, 1
	syscall
	la $a0, result_c # print \n
	li $v0, 4
	syscall

	li $v0, 10 # exit
	syscall


fib:
	addi $sp, $sp, -12
	sw $a0, 8($sp)
	sw $s0, 4($sp)
	sw $ra, 0($sp)

	slti $t0, $a0, 2 # n <= 1 ?
	bne $t0, $zero, return_1

	addi $a0, $a0, -1 # fib(n-1)
	jal fib
	move $s0, $v0

	addi $a0, $a0, -1 # fib(n-2)
	jal fib
	add $v0, $s0, $v0 # fib(n) = fib(n-1) + fib(n-2)
	j fib_return

return_1:
	li $v0, 1

fib_return:
	lw $a0, 8($sp)
	lw $s0, 4($sp)
	lw $ra, 0($sp)
	addi $sp, $sp, 12
	jr $ra

.data
	prompt: .asciiz "Input a number(>= 0): "
	result_a: .asciiz "fib("
	result_b: .asciiz ") = "
	result_c: .byte 0x0a, 0
