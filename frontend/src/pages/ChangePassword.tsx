import {
    Box,
    Button,
    Field,
    Fieldset,
    Flex,
    Heading,
    Icon,
    Input,
    Stack,
    Text,
  } from '@chakra-ui/react';
  import { FaLock } from "react-icons/fa6";  
  export default function ChangePassword() {
    return (
      <Flex
        minH="100vh"
        align="center"
        justify="center"
        bg="gray.50"
        px={4}
      >
        <Box
          bgGradient="linear(to-b, white, gray.50)"
          p={8}
          rounded="lg"
          boxShadow="lg"
          w={{ base: 'full', sm: '400px' }}
          textAlign="center"
        >
          {/* Lock Icon */}
          <Box
            bg="gray.100"
            borderRadius="full"
            w={10}
            h={10}
            mx="auto"
            mb={4}
            display="flex"
            alignItems="center"
            justifyContent="center"
          >
            <Icon as={FaLock} boxSize={5} />
          </Box>
  
          <Heading size="md" mb={1}>
            Set new password
          </Heading>
          <Text fontSize="sm" color="gray.500" mb={6}>
            You can reset your password at any time
          </Text>
  
          <Fieldset.Root size="md" mb={6}>
            <Stack gap={4}>
              <Fieldset.Content>
                {/* Password */}
                <Field.Root>
                  <Field.Label>Password</Field.Label>
                  <Input
                    type="password"
                    name="password"
                    placeholder="Enter your password"
                  />
                </Field.Root>
  
                {/* Confirm Password */}
                <Field.Root>
                  <Field.Label>Confirm Password</Field.Label>
                  <Input
                    type="password"
                    name="confirmPassword"
                    placeholder="********"
                  />
                </Field.Root>
              </Fieldset.Content>
            </Stack>
          </Fieldset.Root>
  
          {/* Reset Button */}
          <Button
            colorScheme="purple"
            w="full"
            fontWeight="medium"
          >
            Reset Password
          </Button>
        </Box>
      </Flex>
    );
  }
  