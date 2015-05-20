.text
.globl main
	main:
	la $a0, query					 #First the query
	li $v0, 4
	syscall
	li $v0, 5					 #Read the input
	syscall
	move $t0, $v0					 #store the value in a temporary variable
#store the base values in $t1, $t2
# $t1 iterates from m-1 to 1
# $t2 maintains a counter of the number of coprimes less than m

	addi $t1, $t0, -1
	li $t2, 0

tot:
	ble $t1, $zero, done  				#termination condition
	move $a0, $t0					#Argument passing
	move $a1, $t1   				#Argument passing 
	jal gcd						#to GCD function
	addi $t3, $v0, -1 
	beqz $t3, inc   				#checking if gcd is one
	addi $t1, $t1, -1				#decrementing the iterator
	j tot 

inc:
	addi $t2, $t2, 1				#incrementing the counter
	addi $t1, $t1, -1				#decrementing the iterator
	j tot

gcd:							#recursive definition
	addi $sp, $sp, -12
	sw $a1, 8($sp)
	sw $a0, 4($sp)
	sw $ra, 0($sp)
	move $v0, $a0
	beqz $a1, gcd_return			        #termination condition
	move $t4, $a0					#computing GCD
	move $a0, $a1
	rem $a1, $t4, $a1
	jal gcd
	lw $a1, 8($sp)
	lw $a0, 4($sp)

gcd_return:
	lw $ra, 0($sp)
	addi $sp, $sp, 12
	jr $ra

done:							 #print the result
#first the message
	la $a0, result_msg
	li $v0, 4
	syscall
#then the value
	move $a0, $t2
	li $v0, 1
	syscall
#exit
	li $v0, 10
	syscall

.data
	query: .asciiz "Input m =  "
	result_msg: .asciiz "Totient(m) =  "
