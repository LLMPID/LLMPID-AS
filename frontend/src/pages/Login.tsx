import {
  Box,
  Button,
  Field,
  Fieldset,
  Flex,
  Group,
  Heading,
  IconButton,
  Input,
  Link,
  Stack,
  Text,
  Image,
} from "@chakra-ui/react";
import { useState } from "react";
import { FiEye, FiEyeOff } from "react-icons/fi";
import { login } from "@/auth/auth"; // path to your login function
import { useSetAtom } from "jotai";
import { authTokenAtom } from "@/atoms/authAtom";
import { toaster, Toaster } from "@/components/ui/toaster"

export default function Login() {
  const [showPassword, setShowPassword] = useState(false);
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [hasError, setHasError] = useState(false);

  const setToken = useSetAtom(authTokenAtom);

  const handleLogin = async () => {
    setLoading(true);
    setHasError(false);
    try {
      const response = await login({ username, password });

      if (response?.access_token) {
        setToken(response.access_token);
      } else {
        setHasError(true);
        toaster.create({
          title: "Login failed",
          description: "Invalid credentials or no token returned.",
          type: "error",
          duration: 4000,
        });
      }
    } catch (error) {
      setHasError(true);
      toaster.create({
        title: "Login error",
        description: "Wrong credentials!",
        type: "error",
        duration: 4000,
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <Flex minH="100vh" align="center" justify="center" bg="gray.50">
      <Box
        bg="white"
        p={8}
        rounded="lg"
        boxShadow="lg"
        w={{ base: "full", sm: "400px" }}
        textAlign="center"
      >
        <Box mb={4}>
          <Image
            src="/logo.png"
            alt="Company Logo"
            mx="auto"
            boxSize="40px"
            objectFit="contain"
          />
        </Box>

        <Heading size="md" mb={1}>
          Welcome 👋
        </Heading>
        <Text fontSize="sm" color="gray.500" mb={6}>
          Please enter your details to sign in
        </Text>

        <Fieldset.Root size="md" mb={4} invalid = {hasError}>
          <Stack gap={4}>
            <Fieldset.Legend srOnly>Sign In</Fieldset.Legend>

            <Fieldset.Content>
              <Field.Root invalid = {hasError}>
                <Field.Label>Username</Field.Label>
                <Input
                  name="username"
                  placeholder="Enter your username"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                />
              </Field.Root>

              <Field.Root invalid = {hasError}>
                <Field.Label>Password</Field.Label>
                <Group attached w="full">
                  <Input
                    name="password"
                    type={showPassword ? "text" : "password"}
                    placeholder="Enter your password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    flex="1"
                  />
                  <IconButton
                    zIndex={-1}
                    aria-label={
                      showPassword ? "Hide password" : "Show password"
                    }
                    variant="outline"
                    onClick={() => setShowPassword(!showPassword)}
                    
                  >
                    {showPassword ? <FiEyeOff /> : <FiEye />}
                  </IconButton>
                </Group>
              </Field.Root>
            </Fieldset.Content>
          </Stack>
        </Fieldset.Root>

        <Flex justify="flex-end" mb={4}>
          <Link fontSize="sm" color="blue.500" href="#">
            Change Password
          </Link>
        </Flex>

        <Button
          colorScheme="purple"
          w="full"
          mb={4}
          onClick={handleLogin}
          loading={loading}
        >
          Sign in
        </Button>

        <Flex justify="center">
          <Link
            href="https://github.com/LLMPID/LLMPID-AS"
            fontSize="sm"
            color="blue.600"
            target="_blank"
            rel="noopener noreferrer"
          >
            📘 View the docs
          </Link>
        </Flex>
      </Box>
      <Toaster/>
    </Flex>
    
    
  );
}
