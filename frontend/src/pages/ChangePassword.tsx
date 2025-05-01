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
  HStack,
} from "@chakra-ui/react";
import { toaster, Toaster } from "@/components/ui/toaster";
import { FaLock, FaArrowLeft } from "react-icons/fa6";
import { Link } from "react-router-dom";
import { useAtomValue } from "jotai";
import { authTokenAtom } from "@/atoms/authAtom";
import { getUsernameFromToken } from "@/utils/jwt";
import { useState } from "react";
import { changePassword } from "@/auth";
import { useSetAtom } from "jotai";
import { useNavigate } from "react-router-dom";

export default function ChangePassword() {
  const token = useAtomValue(authTokenAtom);
  const username = getUsernameFromToken(token);
  const [old_password, setOldPassword] = useState("");
  const [new_password, setNewPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [hasError, setHasError] = useState(false);
  const setToken = useSetAtom(authTokenAtom);
  const navigate = useNavigate();

  const handleChange = async () => {
    setLoading(true);
    setHasError(false);
    try {
      const response = await changePassword({
        username,
        old_password,
        new_password,
      });
      if (response?.access_token) {
        setToken(response.access_token);
        setTimeout(() => {
          navigate("/dashboard");
        }, 2000);
        toaster.create({
          title: "Updated Credentials!",
          description: "You have changed your password",
          type: "success",
          duration: 2000,
        });
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
    <Flex minH="100vh" align="center" justify="center" bg="gray.50" px={4}>
      <Box
        bgGradient="linear(to-b, white, gray.50)"
        p={8}
        rounded="lg"
        boxShadow="lg"
        w={{ base: "full", sm: "400px" }}
        textAlign="center"
      >
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
              <Field.Root invalid={hasError}>
                <Field.Label>Old Password</Field.Label>
                <Input
                  type="password"
                  name="password"
                  placeholder="Enter your password"
                  onChange={(e) => setOldPassword(e.target.value)}
                />
              </Field.Root>

              <Field.Root>
                <Field.Label>New Password</Field.Label>
                <Input
                  type="password"
                  name="confirmPassword"
                  placeholder="********"
                  onChange={(e) => setNewPassword(e.target.value)}
                />
              </Field.Root>
            </Fieldset.Content>
          </Stack>
        </Fieldset.Root>

        <Button
          w="full"
          fontWeight="medium"
          onClick={handleChange}
          loading={loading}
        >
          Reset Password
        </Button>
        <Box mt={4} textAlign="center">
          <Link to="/dashboard">
            <HStack
              justify="center"
              gap={1}
              color="gray.600"
              _hover={{ color: "blue.500" }}
            >
              <Icon as={FaArrowLeft} />
              <Text fontSize="sm">Back to dashboard</Text>
            </HStack>
          </Link>
        </Box>
      </Box>
      <Toaster />
    </Flex>
  );
}
