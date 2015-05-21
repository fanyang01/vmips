.text
.globl main
main:
	li $v0, 4
	la $a0, hint
	syscall

	li $v0, 4
	la $a0, prompt_1
	syscall
	li $v0, 5
	syscall
	move $s0, $v0

	li $v0, 4
	la $a0, prompt_2
	syscall
	li $v0, 5
	syscall
	move $s1, $v0

	li $v0, 4 # print "ack(m, n) = "
	la $a0, result
	syscall

	move $a0, $s0
	move $a1, $s1
	jal ack
	move $a0, $v0 # print result
	li $v0, 1
	syscall

	la $a0, newline # print \n
	li $v0, 4
	syscall

	li $v0, 10 # exit
	syscall

ack:
	addi $sp, $sp, -12
	sw $a0, 8($sp)
	sw $a1, 4($sp)
	sw $ra, 0($sp)

	beq $a0, $zero, case_m_0 # m == 0
	beq $a1, $zero, case_n_0 # n == 0

	addi $a1, $a1, -1 # A(m, n-1)
	jal ack
	move $a1, $v0 # A(m-1, A(m, n-1))
	addi $a0, $a0, -1
	jal ack
	j ack_return

case_m_0:
	addi $v0, $a1, 1 # n + 1
	j ack_return

case_n_0:
	addi $a0, $a0, -1 # A(m-1, 1)
	addi $a1, $a1, 1
	jal ack

ack_return:
	lw $a0, 8($sp)
	lw $a1, 4($sp)
	lw $ra, 0($sp)
	addi $sp, $sp, 12
	jr $ra

.data
hint: .ascii "Ackermann function A(m, n) = :"
	.byte 0x0a
	.ascii "n + 1, if m == 0"
	.byte 0x0a
	.ascii "A(m-1, 1), if m > 0 and n == 0"
	.byte 0x0a
	.ascii "A(m-1, A(m, n-1)), if m > 0 and n > 0"
newline:
	.byte 0x0a, 0
prompt_1: .asciiz "Input m: "
prompt_2: .asciiz "Input n: "
result: .asciiz "A(m, n) = "
